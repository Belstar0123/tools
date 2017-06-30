package main

import (
    "os"
    "flag"
    "io/ioutil"
    "fmt"
    "sort"
    "strings"
    "strconv"

    "github.com/disintegration/imaging"
    "github.com/Belstar0123/tools/lib"
)

type FileInfos []os.FileInfo

// ファイルソート設定
type ByName struct{ FileInfos }

func (fi ByName) Len() int {
    return len(fi.FileInfos)
}

func (fi ByName) Swap(i, j int) {
    fi.FileInfos[i], fi.FileInfos[j] = fi.FileInfos[j], fi.FileInfos[i]
}

func (fi ByName) Less(i, j int) bool {
    return fi.FileInfos[j].Name() > fi.FileInfos[i].Name()
}

/*
Komifloにログインして、漫画を画像ファイルとして保存する.
 */
func main() {
    var argF, argH, argW, argP string
    var trim_height, trim_width int = 0, 0

    // 各オプション
    flag.StringVar(&argF, "f", "", "TargetDirectory(FullPath)")
    flag.StringVar(&argH, "h", "", "Trimming height(pixel)")
    flag.StringVar(&argW, "w", "", "Trimming width(pixel)")
    flag.StringVar(&argP, "p", "", "Trimming base point(Left/Center/Right/Top/TopLeft/TopRight/Bottom/BottomLeft/BottomRight) default Center")

    // コマンドライン引数を解析
    flag.Parse()

    // f引数がない場合はカレントディレクトリを使用
    if argF == "" {
        argF = lib.GetDirName()
    }

    // h引数がある場合は数値変換
    if argH != "" {
        trim_height, _ = strconv.Atoi(argH)
    }

    // w引数がある場合は数値変換
    if argW != "" {
        trim_width, _ = strconv.Atoi(argW)
    }

    // pオプションが無い場合はCenter指定
    var trim_pos_str string = argP
    var trim_pos imaging.Anchor
    switch argP {
    case "Left":
        trim_pos = imaging.Left
    case "Center":
        trim_pos = imaging.Center
    case "Right":
        trim_pos = imaging.Right
    case "Top":
        trim_pos = imaging.Top
    case "TopLeft":
        trim_pos = imaging.TopLeft
    case "TopRight":
        trim_pos = imaging.TopRight
    case "Bottom":
        trim_pos = imaging.Bottom
    case "BottomLeft":
        trim_pos = imaging.BottomLeft
    case "BottomRight":
        trim_pos = imaging.BottomRight
    default:
        trim_pos = imaging.Center
        trim_pos_str = "Center"
    }

    // ディレクトリとファイル名に分割して格納
    var targetPath string = argF

    // ディレクトリが無い場合は、カレントディレクトリを使用
    if targetPath == "" {
        targetPath = lib.GetDirName()
    }

    fmt.Printf("変換対象フォルダ: %s\n", targetPath)
    fmt.Printf("トリミング高: %d\n", trim_height)
    fmt.Printf("トリミング幅: %d\n", trim_width)
    fmt.Printf("トリミング基準: %s\n", trim_pos_str)

    // 取得しようとしているパスがディレクトリかチェック
    var isDir, _ = lib.IsDirectory(argF)

    // ディレクトリならば、そのディレクトリ配下のファイルを調べる。
    if isDir == true {
        targetPath = argF
    }

    // ディレクトリ内のファイル情報の読み込み[] *os.FileInfoが返る。
    fileInfos, err := ioutil.ReadDir(targetPath)

    // ディレクトリの読み込みに失敗したらエラーで終了
    if err != nil {
        fmt.Errorf("Directory cannot read %s\n", err)
        os.Exit(1)
    }

    // 直近のディレクトリ名のみ抽出
    var temps []string = strings.Split(targetPath, "\\")
    var dirName string = temps[len(temps)-1]

    // 作成日付順にソート
    sort.Sort(ByName{fileInfos})
    // ファイルを一つずつ処理する
    var counter int = 1
    for _, fileInfo := range fileInfos {
        // *FileInfo型
        var findName = (fileInfo).Name()
        var extName = lib.GetExtension(findName)
        // ファイルの先頭が"."の場合はスキップ
        unicodeName := []rune(findName)
        compCode := []rune(".")
        if unicodeName[0] == compCode[0] {
            continue
        }
        // トリミング実施
        lib.ExecTrimming(trim_width, trim_height, trim_pos, targetPath + findName)
        // リネーム後のファイル名
        rename := dirName + "_" + fmt.Sprintf("%04d", counter) + extName
        // ファイル名をリネーム
        if err := os.Rename(targetPath+"\\"+findName, targetPath+"\\"+rename); err != nil {
            fmt.Println(err)
        }
        counter++
    }

}
