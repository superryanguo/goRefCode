package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

type GMCrypt struct {
	PublicFile  string
	PrivateFile string
}

var (
	path        = "./"
	privateFile = "sm2private.pem" // 私钥文件
	publicFile  = "sm2public.pem"  // 公钥文件
	data        = "hello 国密"
)

// 测试一下
func main() {
	GenerateSM2Key() // 密钥生成并保存在文件中
	crypt := GMCrypt{
		PublicFile:  path + publicFile,
		PrivateFile: path + privateFile,
	}
	encryptText, _ := crypt.Encrypt(data) // 加密
	fmt.Println(encryptText)
	decryptText, _ := crypt.Decrypt(encryptText) // 解密
	fmt.Println(decryptText)

	msg := []byte("hello 国密")
	sig, key, _ := CreateSm2Sig(msg) // 签名
	fmt.Printf("签名结果：%x\n公钥：%v, \n", sig, key)
	verSm2Sig := VerSm2Sig(key, msg, sig) // 验证签名
	fmt.Println("验证结果为：", verSm2Sig)
}

// 生成公钥、私钥
func GenerateSM2Key() {
	// 生成私钥、公钥
	priKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("秘钥产生失败：", err)
		os.Exit(1)
	}
	pubKey := &priKey.PublicKey
	// 生成文件 保存私钥、公钥
	// x509 编码
	pemPrivKey, _ := x509.WritePrivateKeyToPem(priKey, nil)
	privateFile, _ := os.Create(path + privateFile)
	defer privateFile.Close()
	privateFile.Write(pemPrivKey)
	pemPublicKey, _ := x509.WritePublicKeyToPem(pubKey)
	publicFile, _ := os.Create(path + publicFile)
	defer publicFile.Close()
	publicFile.Write(pemPublicKey)
}

// 读取密钥文件
func readPemCxt(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return []byte{}, err
	}
	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	if err != nil {
		return []byte{}, err
	}
	return buf, err
}

// 加密
func (s *GMCrypt) Encrypt(data string) (string, error) {
	pub, err := readPemCxt(s.PublicFile)
	if err != nil {
		return "", err
	}
	// read public key
	publicKeyFromPem, err := x509.ReadPublicKeyFromPem(pub)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	ciphertxt, err := publicKeyFromPem.EncryptAsn1([]byte(data), rand.Reader)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertxt), nil
}

// 解密
func (s *GMCrypt) Decrypt(data string) (string, error) {
	pri, err := readPemCxt(s.PrivateFile)
	if err != nil {
		return "", err
	}
	privateKeyFromPem, err := x509.ReadPrivateKeyFromPem(pri, nil)
	if err != nil {
		return "", err
	}
	ciphertxt, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	plaintxt, err := privateKeyFromPem.DecryptAsn1(ciphertxt)
	if err != nil {
		return "", err
	}
	return string(plaintxt), nil
}

// 使用私钥创建签名
func CreateSm2Sig(msg []byte) ([]byte, *sm2.PublicKey, error) {
	// 读取密钥
	pri, _ := readPemCxt(path + privateFile)
	privateKey, _ := x509.ReadPrivateKeyFromPem(pri, nil)
	c := sm2.P256Sm2() // 椭圆曲线
	priv := new(sm2.PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = privateKey.D
	priv.PublicKey.X = privateKey.X
	priv.PublicKey.Y = privateKey.Y
	sign, err := priv.Sign(rand.Reader, msg, nil) // sm2签名
	if err != nil {
		return nil, nil, err
	}
	return sign, &priv.PublicKey, err
}

// 验证签名
func VerSm2Sig(pub *sm2.PublicKey, msg []byte, sign []byte) bool {
	isok := pub.Verify(msg, sign)
	return isok
}
