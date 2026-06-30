package main

import (
	"fmt"
	"net/http"
)

// 前回の構造体（簡略化しています）
type TrainInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func delayHandler(w http.ResponseWriter, r *http.Request) {
	// 本来はここで「駅すぱあとAPI」からデータを取ってくる
	infoList := []TrainInfo{
		{Name: "JR京都線", Status: "平常運転"},
		{Name: "御堂筋線", Status: "運転見合わせ"},
	}

	// HTMLとしてブラウザに返す
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<h1>🚋 現在の運行情報 </h1>")
	for _, info := range infoList {
		fmt.Fprintf(w, "<p>【%s】: %s</p>", info.Name, info.Status)
	}
}

func main() {
	// 「/delay」というURLにアクセスされたら delayHandler を実行する
	http.HandleFunc("/delay", delayHandler)

	fmt.Println("サーバーを起動しました。ブラウザで http://localhost:8080/delay を開いてください。")
	// 8080番ポートで待ち受け開始
	http.ListenAndServe(":8080", nil)
}