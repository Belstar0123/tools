package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Belstar0123/tools/lib"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type EventData struct {
	Name                   string
	Id                     string
	MaxNo                  int
	AddPrefix              string
	BaseFormat             string
	Mp3LinesCpMaxNo        int
	Mp3LinesQMaxNo         int
	Mp3LinesSMaxNo         int
	Mp3LinesSSubMaxNo      int
	Mp3LinesSPrefix        []string
	Mp3NaviCharacterId     []string
	Mp3NaviCharacterMaxNo  int
	JpgNaviCharacterId     []string
	JpgNaviCharacterMaxNo  int
	JpgNaviCharacterPrefix []string
}

type BasicData struct {
	Name string
	Id   string
}

func readConfig() (err error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("設定ファイル読み込みエラー: %s\n", err)
	}
	return
}

func arrayContains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func main() {
	// コマンドライン オプションの定義
	var argB, argE, argR, argSR, argSSR, argN, argS, argP string
	flag.StringVar(&argB, "b", "", "ON: target output")
	flag.StringVar(&argE, "e", "", "target output event id")
	flag.StringVar(&argR, "r", "", "target output character id")
	flag.StringVar(&argSR, "sr", "", "target output character id")
	flag.StringVar(&argSSR, "ssr", "", "target output character id")
	flag.StringVar(&argN, "n", "", "ON: target output")
	flag.StringVar(&argS, "s", "", "ON: target output")
	flag.StringVar(&argP, "p", "", "ON: target output")

	// コマンドライン引数を解析
	flag.Parse()

	// 設定ファイルの読み込み
	err := readConfig()
	if err != nil {
		// 設定ファイルが読めなかったら処理終了
		return
	}

	// 動作ログの取得
	operation_result_path := viper.GetString("output.operation_result_file_path")
	fmt.Printf("動作ログ出力先: %s\n", operation_result_path)
	var r_list, sr_list, ssr_list, e_list []string
	fp, err := os.Open(operation_result_path)
	defer fp.Close()
	if err != nil {
		// 空ファイルを作成して読み直す
		ioutil.WriteFile(operation_result_path, []byte(""), os.ModePerm)
		fp, err = os.Open(operation_result_path)
		if err != nil {
			// 作成したファイルが読めなかったから処理終了
			return
		}
	}
	reader := bufio.NewReaderSize(fp, 51200)
	line_count := 0
	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		line_count++
		switch line_count {
		case 1:
			r_list = strings.Split(line, ",")
		case 2:
			sr_list = strings.Split(line, ",")
		case 3:
			ssr_list = strings.Split(line, ",")
		case 4:
			e_list = strings.Split(line, ",")
		default:
		}
	}

	// 出力先の取得
	output_path := viper.GetString("output.path")
	fmt.Printf("出力先: %s\n", output_path)

	// 共通背景処理
	if argB == "ON" {
		common_bg_base_url := viper.GetString("background.common.base_url")
		common_bg_max_no := viper.GetInt("background.common.max_no")
		if common_bg_base_url != "" {
			fmt.Println("共通背景出力処理開始")
			// 保存先のディレクトリ設定
			common_bg_output_path := output_path + "/背景/共通/"
			lib.CreateDirectory(common_bg_output_path)
			// 共通背景データの保存
			for no := 0; no <= common_bg_max_no; no++ {
				target_url := fmt.Sprintf(common_bg_base_url, no)
				filename := lib.GetFilename(target_url)
				target_file_path := common_bg_output_path + filename
				lib.DownloadToFile(target_url, target_file_path)
			}
			fmt.Println("  └ 共通背景出力処理終了")
		} else {
			fmt.Errorf("共通背景設定なし\n")
		}

		// イベント背景処理
		event_bg_base_url := viper.GetString("background.event.base_url")
		event_bg_max_no := viper.GetInt("background.event.max_no")
		if event_bg_base_url != "" {
			fmt.Println("イベント背景出力処理開始")
			// 保存先のディレクトリ設定
			event_bg_output_path := output_path + "/背景/イベント/"
			lib.CreateDirectory(event_bg_output_path)
			// 共通背景データの保存
			for no := 0; no <= event_bg_max_no; no++ {
				target_url := fmt.Sprintf(event_bg_base_url, no)
				filename := lib.GetFilename(target_url)
				target_file_path := event_bg_output_path + filename
				lib.DownloadToFile(target_url, target_file_path)
			}
			fmt.Println("  └ イベント背景出力処理終了")
		} else {
			fmt.Errorf("イベント背景設定なし\n")
		}
	}

	// イベント画像処理
	if argE != "" {
		event_base_url := viper.GetString("event.base_url")
		event_base_mp3_url := viper.GetString("event.base_mp3_url")
		event_base_navi_url := viper.GetString("event.base_navi_url")
		event_base_navi_jpg_url := viper.GetString("event.base_navi_jpg_url")
		event_lists := viper.Get("event.lists")
		if event_base_url != "" {
			fmt.Println("イベント画像出力処理開始")
			// yamlデータを配列に格納
			var event_data []EventData
			for _, v := range event_lists.([]interface{}) {
				var new_event_data EventData = EventData{}
				new_event_data.Mp3LinesSPrefix = make([]string, 0)
				new_event_data.Mp3NaviCharacterId = make([]string, 0)
				for sk, sv := range v.(map[interface{}]interface{}) {
					switch sk {
					case "name":
						if sv != nil {
							new_event_data.Name = sv.(string)
						}
					case "id":
						if sv != nil {
							new_event_data.Id = sv.(string)
						}
					case "max_no":
						if sv != nil {
							new_event_data.MaxNo = sv.(int)
						}
					case "add_prefix":
						if sv != nil {
							new_event_data.AddPrefix = sv.(string)
						}
					case "base_format":
						if sv != nil {
							new_event_data.BaseFormat = sv.(string)
						}
					case "mp3_lines_cp_max_no":
						if sv != nil {
							new_event_data.Mp3LinesCpMaxNo = sv.(int)
						}
					case "mp3_lines_q_max_no":
						if sv != nil {
							new_event_data.Mp3LinesQMaxNo = sv.(int)
						}
					case "mp3_lines_s_max_no":
						if sv != nil {
							new_event_data.Mp3LinesSMaxNo = sv.(int)
						}
					case "mp3_lines_s_sub_max_no":
						if sv != nil {
							new_event_data.Mp3LinesSSubMaxNo = sv.(int)
						}
					case "mp3_lines_s_prefix":
						if sv != nil {
							for _, v := range sv.([]interface{}) {
								new_event_data.Mp3LinesSPrefix = append(new_event_data.Mp3LinesSPrefix, v.(string))
							}
						}
					case "mp3_navi_character_id":
						if sv != nil {
							for _, v := range sv.([]interface{}) {
								new_event_data.Mp3NaviCharacterId = append(new_event_data.Mp3NaviCharacterId, v.(string))
							}
						}
					case "mp3_navi_character_max_no":
						if sv != nil {
							new_event_data.Mp3NaviCharacterMaxNo = sv.(int)
						}
					case "jpg_navi_character_id":
						if sv != nil {
							for _, v := range sv.([]interface{}) {
								new_event_data.JpgNaviCharacterId = append(new_event_data.JpgNaviCharacterId, v.(string))
							}
						}
					case "jpg_navi_character_max_no":
						if sv != nil {
							new_event_data.JpgNaviCharacterMaxNo = sv.(int)
						}
					case "jpg_navi_character_prefix":
						if sv != nil {
							for _, v := range sv.([]interface{}) {
								new_event_data.JpgNaviCharacterPrefix = append(new_event_data.JpgNaviCharacterPrefix, v.(string))
							}
						}
					}
				}
				if argE != new_event_data.Id {
					continue
				}
				event_data = append(event_data, new_event_data)
			}

			// 配列データを処理
			for _, v := range event_data {
				if !arrayContains(e_list, v.Id) {
					fmt.Printf("  ├ イベント[%s] 処理中\n", v.Name)
					event_output_path := output_path + "/イベント/" + v.Name + "/"
					lib.CreateDirectory(event_output_path)
					for n := 0; n <= v.MaxNo; n++ {
						target_url := fmt.Sprintf(event_base_url, v.Id, v.BaseFormat)
						target_url = fmt.Sprintf(target_url, n)
						filename := lib.GetFilename(target_url)
						target_file_path := event_output_path + filename
						lib.DownloadToFile(target_url, target_file_path)
					}

					// prefix分のループ
					prefixs := strings.Split(v.AddPrefix, ",")
					for _, prefix := range prefixs {
						for n := 0; n <= v.MaxNo; n++ {
							target_url := fmt.Sprintf(event_base_url, v.Id+"_"+prefix, v.BaseFormat)
							target_url = fmt.Sprintf(target_url, n)
							filename := lib.GetFilename(target_url)
							target_file_path := event_output_path + filename
							lib.DownloadToFile(target_url, target_file_path)
						}
					}

					// base_mp3データ
					for n := 0; n <= v.Mp3LinesCpMaxNo; n++ {
						for m := 0; m <= v.Mp3LinesQMaxNo; m++ {
							for k := 0; k <= v.Mp3LinesSMaxNo; k++ {
								for i := 0; i <= v.Mp3LinesSSubMaxNo; i++ {
									for _, vv := range v.Mp3LinesSPrefix {
										target_url := fmt.Sprintf(event_base_mp3_url, v.Id, n, m, k, i, vv)
										filename := lib.GetFilename(target_url)
										target_file_path := event_output_path + filename
										lib.DownloadToFile(target_url, target_file_path)
									}
								}
							}
						}
					}

					// base_navi_mp3データ
					for _, m := range v.Mp3NaviCharacterId {
						for k := 0; k <= v.Mp3NaviCharacterMaxNo; k++ {
							target_url := fmt.Sprintf(event_base_navi_url, m, k)
							filename := lib.GetFilename(target_url)
							target_file_path := event_output_path + filename
							lib.DownloadToFile(target_url, target_file_path)
						}
					}

					// base_navi_jpgデータ
					for n := 0; n <= v.JpgNaviCharacterMaxNo; n++ {
						for _, m := range v.JpgNaviCharacterPrefix {
							for _, k := range v.JpgNaviCharacterId {
								target_url := ""
								if m != "" {
									target_url = fmt.Sprintf(event_base_navi_jpg_url, k, m+"_", n)
								} else {
									target_url = fmt.Sprintf(event_base_navi_jpg_url, k, m, n)
								}
								filename := lib.GetFilename(target_url)
								target_file_path := event_output_path + filename
								lib.DownloadToFile(target_url, target_file_path)
							}
						}
					}
					e_list = append(e_list, v.Id)
				} else {
					fmt.Println("  ├ イベント画像出力スキップ")
				}
			}
			fmt.Println("  └ イベント画像出力処理終了")
		} else {
			fmt.Errorf("イベント画像設定なし\n")
		}
	}

	// Rキャラクター画像の処理
	if argR != "" {
		r_character_members := viper.Get("character.R.members")
		getCharacterGraphic("R", r_character_members, output_path, argR, r_list)
	}
	// SRキャラクター画像の処理
	if argSR != "" {
		sr_character_members := viper.Get("character.SR.members")
		getCharacterGraphic("SR", sr_character_members, output_path, argSR, sr_list)
	}
	// SSRキャラクター画像の処理
	if argSSR != "" {
		ssr_character_members := viper.Get("character.SSR.members")
		getCharacterGraphic("SSR", ssr_character_members, output_path, argSSR, ssr_list)
	}
	// NPCキャラクター画像の処理
	if argN == "ON" {
		sp_character_members := viper.Get("character.sp.members")
		getCharacterGraphic("sp", sp_character_members, output_path, "", []string{})
		// その他キャラクター画像の処理
		mob_character_members := viper.Get("character.npc.members")
		getCharacterGraphic("npc", mob_character_members, output_path, "", []string{})
	}

	// 召喚獣画像の処理
	if argS == "ON" {
		n_summonstone_members := viper.Get("summon.RN.members")
		getSummonStoneGraphic("N", n_summonstone_members, output_path)
		r_summonstone_members := viper.Get("summon.R.members")
		getSummonStoneGraphic("R", r_summonstone_members, output_path)
		sr_summonstone_members := viper.Get("summon.SR.members")
		getSummonStoneGraphic("SR", sr_summonstone_members, output_path)
		ssr_summonstone_members := viper.Get("summon.SSR.members")
		getSummonStoneGraphic("SSR", ssr_summonstone_members, output_path)

	}

	// PC画像の処理
	if argP == "ON" {
		pc_base_url := viper.GetString("pc.base_url")
		pc_base_sexs := viper.Get("pc.base_sexs")
		pc_base_jobs := viper.Get("pc.base_jobs")

		// キャラクター画像処理
		if pc_base_url != "" {
			fmt.Print("PCジョブ画像出力処理開始\n")
			// yamlデータを配列に格納
			var pc_sex_data []BasicData
			for _, v := range pc_base_sexs.([]interface{}) {
				new_pc_sex_data := BasicData{}
				for sk, sv := range v.(map[interface{}]interface{}) {
					switch sk {
					case "name":
						new_pc_sex_data.Name = sv.(string)
					case "id":
						new_pc_sex_data.Id = strconv.Itoa(sv.(int))
					}
				}
				pc_sex_data = append(pc_sex_data, new_pc_sex_data)
			}
			var pc_job_data []BasicData
			for _, v := range pc_base_jobs.([]interface{}) {
				new_pc_job_data := BasicData{}
				for sk, sv := range v.(map[interface{}]interface{}) {
					switch sk {
					case "name":
						new_pc_job_data.Name = sv.(string)
					case "id":
						new_pc_job_data.Id = sv.(string)
					}
				}
				pc_job_data = append(pc_job_data, new_pc_job_data)
			}

			// 配列データを処理
			for _, v := range pc_sex_data {
				fmt.Printf("処理中... %s\n", v.Name)
				pc_output_path := output_path + "/プレイヤー/" + v.Name
				lib.CreateDirectory(pc_output_path)
				// ジョブ画像
				for _, job := range pc_job_data {
					target_url := fmt.Sprintf(pc_base_url, job.Id, v.Id)
					filename := lib.GetFilename(target_url)
					target_file_path := pc_output_path + "/" + job.Name + "_" + filename
					lib.DownloadToFile(target_url, target_file_path)
				}
			}
			fmt.Print("  └ PCジョブ画像出力処理終了\n")
		} else {
			fmt.Errorf("PCジョブ画像設定なし\n")
		}
	}

	// 処理済みリストの書き出し
	// Rキャラリスト
	ioutil.WriteFile(operation_result_path, []byte(strings.Join(r_list, ",")), os.ModePerm)
	// SRキャラリスト
	ioutil.WriteFile(operation_result_path, []byte(strings.Join(sr_list, ",")), os.ModePerm)
	// SSRキャラリスト
	ioutil.WriteFile(operation_result_path, []byte(strings.Join(ssr_list, ",")), os.ModePerm)
	// イベントリスト
	ioutil.WriteFile(operation_result_path, []byte(strings.Join(ssr_list, ",")), os.ModePerm)

}

// レアリティ毎のキャラ画像データを保存する
func getCharacterGraphic(rarity string, character_members interface{}, output_path string, targetId string, exclusion_list []string) {
	if character_members == nil {
		return
	}

	// キャラクター画像の処理
	character_base_url := viper.GetString("character.base_url")
	character_base_prefixs := viper.Get("character.base_prefix")
	character_add_url := viper.GetString("character.add_url")
	character_add_prefix_bases := viper.Get("character.add_prefix_base")
	character_add_prefixs := viper.Get("character.add_prefix")
	character_add_prefix_nos := viper.Get("character.add_prefix_no")

	// キャラクター画像処理
	if character_base_url != "" {
		fmt.Printf("%sキャラクター画像出力処理開始\n", rarity)
		// yamlデータを配列に格納
		var character_data []BasicData
		for _, v := range character_members.([]interface{}) {
			new_character_data := BasicData{}
			for sk, sv := range v.(map[interface{}]interface{}) {
				switch sk {
				case "name":
					new_character_data.Name = sv.(string)
				case "id":
					new_character_data.Id = strconv.Itoa(sv.(int))
				}
			}
			if targetId != "" {
				if targetId != "ON" && targetId != new_character_data.Id {
					continue
				}
			}
			character_data = append(character_data, new_character_data)
		}

		// 配列データを処理
		for _, v := range character_data {
			fmt.Printf("処理中... %s\n", v.Name)
			if (!arrayContains(exclusion_list, v.Id)) {
				character_output_path := output_path + "/キャラクター/" + v.Name
				lib.CreateDirectory(character_output_path)
				// キャラ画像(開放絵)
				for _, base_prefix := range character_base_prefixs.([]interface{}) {
					target_url := fmt.Sprintf(character_base_url, v.Id, base_prefix)
					filename := lib.GetFilename(target_url)
					target_file_path := character_output_path + "/" + rarity + "_" + filename
					lib.DownloadToFile(target_url, target_file_path)
				}
				// 表情
				for _, base_prefix := range character_add_prefix_bases.([]interface{}) {
					target_prefix := v.Id
					if len(base_prefix.(string)) > 0 {
						target_prefix = target_prefix + "_" + base_prefix.(string)
					}
					for _, add_prefix := range character_add_prefixs.([]interface{}) {
						target_prefix_add := target_prefix
						if len(add_prefix.(string)) > 0 {
							target_prefix_add = target_prefix_add + "_" + add_prefix.(string)
						}
						target_url := fmt.Sprintf(character_add_url, target_prefix_add)
						filename := lib.GetFilename(target_url)
						target_file_path_add := character_output_path + "/" + rarity + "_" + filename
						lib.DownloadToFile(target_url, target_file_path_add)

						// 表情差分
						for _, no := range character_add_prefix_nos.([]interface{}) {
							target_prefix_diff := target_prefix_add + no.(string)
							target_url_diff := fmt.Sprintf(character_add_url, target_prefix_diff)
							filename_diff := lib.GetFilename(target_url_diff)
							target_file_path_diff := character_output_path + "/" + rarity + "_" + filename_diff
							lib.DownloadToFile(target_url_diff, target_file_path_diff)
						}
					}
				}
				exclusion_list = append(exclusion_list, v.Id)
			} else {
				fmt.Printf("  ├ %sキャラクター画像出力スキップ\n", rarity)
			}
		}
		fmt.Printf("  └ %sキャラクター画像出力処理終了\n", rarity)
	} else {
		fmt.Errorf("%sキャラクター画像設定なし\n", rarity)
	}
}

// レアリティ毎の召喚石画像データを保存する
func getSummonStoneGraphic(rarity string, summonstone_members interface{}, output_path string) {
	if summonstone_members == nil {
		return
	}

	// 召喚石画像の処理
	summonstone_base_url := viper.GetString("summon.base_url")

	// キャラクター画像処理
	if summonstone_base_url != "" {
		fmt.Printf("%s召喚石画像出力処理開始\n", rarity)
		// yamlデータを配列に格納
		var summonstone_data []BasicData
		for _, v := range summonstone_members.([]interface{}) {
			new_summonstone_data := BasicData{}
			for sk, sv := range v.(map[interface{}]interface{}) {
				switch sk {
				case "name":
					new_summonstone_data.Name = sv.(string)
				case "id":
					new_summonstone_data.Id = strconv.Itoa(sv.(int))
				}
			}
			summonstone_data = append(summonstone_data, new_summonstone_data)
		}

		// 配列データを処理
		for _, v := range summonstone_data {
			summonstone_output_path := output_path + "/召喚石"
			lib.CreateDirectory(summonstone_output_path)
			// 召喚石画像
			target_url := fmt.Sprintf(summonstone_base_url, v.Id)
			filename := lib.GetFilename(target_url)
			target_file_path := summonstone_output_path + "/" + rarity + "_" + v.Name + "_" + filename
			lib.DownloadToFile(target_url, target_file_path)
		}
		fmt.Printf("  └ %s召喚石画像出力処理終了\n", rarity)
	} else {
		fmt.Errorf("%s召喚石画像設定なし\n", rarity)
	}
}
