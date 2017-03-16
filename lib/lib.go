package lib

import (
	"time"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/disintegration/imaging"
	"image"
)

/*
Parameter
  s 待機秒数
  c 結果引き渡し用channel
*/
func SleepSec(s time.Duration, c chan string) {
	log.Printf("%dsec sleep1 started.", s)
	time.Sleep(s * time.Second)
	if c != nil {
		c <- fmt.Sprintf("%dsec wait finished.", s)
	} else {
		log.Printf("%dsec sleep1 finished.", s)
	}
}

/*
指定されたファイル名がディレクトリかどうか調べる
*/
func IsDirectory(name string) (isDir bool, err error) {
	fInfo, err := os.Stat(name) // FileInfo型が返る。
	if err != nil {
		isDir = false
		return
	}
	// ディレクトリかどうかチェック
	isDir = fInfo.IsDir()
	return
}

/*
指定されたファイル名から拡張子を抽出
 */
func GetExtension(fileName string) (ext string) {
	pos := strings.LastIndex(fileName, ".")
	ext = fileName[pos:]
	return
}

/*
カレントディレクトリの取得
*/
func GetDirName() (ret string){
	var curDir, _ = os.Getwd()
	return curDir + "/"
}

/*
指定サイズに画像をトリミング

Parameter
  w トリミング幅
  h トリミング高
  p トリミング基準点
  path 対象画像フルパス
 */
func ExecTrimming(w int, h int, p imaging.Anchor, path string) {
	// トリミング幅高が指定されていたらトリミングを実施する
	if h != 0 && w != 0 {
		// 画像をトリミングする
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("target_file: %s\n", path)
			return
		}
		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("target_file: %s\n", path)
			return
		}
		dstImage := imaging.Fill(img, w, h, p, imaging.Lanczos)
		// ファイル保存
		err = imaging.Save(dstImage, path)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("target_file: %s\n", path)
			return
		}
	}
}