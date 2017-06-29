package lib

import (
	"time"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/disintegration/imaging"
	"image"
	"net/http"
	"io"
)

type DownloadError struct {
	StatusCode int
}
func (e *DownloadError) Error() string {
	return "httpd error!"
}

type SaveError struct {
	Filename, Message string
}
func (e *SaveError) Error() string {
	return e.Message
}

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

/*
指定されたディレクトリを作成
 */
func CreateDirectory(path string) (err error) {
	if err = os.MkdirAll(path, 0777); err != nil {
		fmt.Errorf("ディレクトリの作成に失敗しました: %s\n", err)
	}
	return
}

/*
指定されたローカルファイルの存在チェック
 */
func IsExists(path string) (isExist bool) {
	_, err := os.Stat(path) // FileInfo型が返る。
	isExist = !os.IsNotExist(err)
	return
}

/*
指定されたURLからファイル名を取得
 */
func GetFilename(url string) (filename string) {
	pos := strings.LastIndex(url, "/") + 1
	filename = url[pos:]
	return
}

/*
指定されたURLの画像を指定されたローカルファイルに保存
 */
func DownloadToFile(url string, path string) (err error) {
	// 既にローカルファイルが存在するかどうかチェック
	if IsExists(path) {
		// 存在していたら処理終了
		return nil
	}

	// URLのデータを取得
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Errorf("URLからのデータ取得に失敗しました :%s\n", err)
		return err
	}
	// 200以外のレスポンスはダウンロードエラーとする
	if response.StatusCode != http.StatusOK {
		return &DownloadError{StatusCode: response.StatusCode}
	}
	// 書き込むファイルの準備
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		fmt.Errorf("ファイルの書き込みに失敗しました :%s\n", err)
		return err
	}

	// データをファイルに書き込み
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return &SaveError{Message: err.Error(), Filename: file.Name()}
	}
	return nil
}