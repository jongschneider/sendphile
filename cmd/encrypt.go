/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var key *string
var file *string

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypts a file using AES",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Println("file:", *file)
		fmt.Println("key:", *key)
		data, err := ioutil.ReadFile(*file)
		if err != nil {
			return err
		}

		b := encrypt(data, *key)

		out, err := createDstFilepath(*file)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(out, b, 0644)
		if err != nil {
			return err
		}

		fmt.Printf("contents of %s were encrypted and stored in %s\n", *file, out)
		return nil
	},
}

func init() {
	// rootCmd.AddCommand(encryptCmd)
	file = encryptCmd.Flags().StringP("file", "f", "", "the filename to encrypt")
	key = encryptCmd.Flags().StringP("key", "k", "", "the key used to encrypt and decrypt the file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// The following encryption process comes from https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

// createHash takes a passphrase or any string, hashes it, then returns the hash as a hexadecimal value.
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

// createDstFilepath creates a unique filename for the destination file which will store our encrypted data.
func createDstFilepath(in string) (string, error) {
	dir, filename := path.Split(in)

	// Retrieve all files in the directory so we can diff the new filename against the existing filenames in the dir. We want to avoid naming collisions.
	files, err := getFiles(dir)
	if err != nil {
		return "", err
	}

	// Create the initial version of the new filename. If this name is not already in the dir, it will become the name of the file. If this name already exists, we will append a number on the end to distinguish this file from others of the same name.
	out := fmt.Sprintf("enc_%s", filename)

	if _, exists := files[out]; exists {
		for count := 1; ; count++ {
			version := fmt.Sprintf("(%d)", count)
			base := strings.Split(out, path.Ext(filename))[0]
			filename = fmt.Sprintf("%s%s%s", base, version, path.Ext(filename))
			if _, exists := files[filename]; !exists {
				out = filename
				break
			}
		}
	}

	out = path.Join(dir, out)
	return out, nil
}

// getFiles stores a list of files in a lookup map.
func getFiles(dir string) (files map[string]struct{}, err error) {
	files = make(map[string]struct{})

	if dir == "" {
		dir = "./"
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files[path] = struct{}{}
		return nil
	})

	return
}
