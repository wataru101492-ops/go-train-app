package main

import (
	"fmt"
	"net/http"
	"os"
)

type TrainInfo struct {
	Name    string `json:"name"`
	Company string `json:"company"` // 鉄道会社名も追加
	Status  string `json:"status"`
}

func delayHandler(w http.ResponseWriter, r *http.Request) {
	// 表示確認用のデータ
	infoList := []TrainInfo{
		{Name: "JR京都線", Company: "JR西日本", Status: "平常運転"},
		{Name: "御堂筋線", Company: "大阪メトロ", Status: "運転見合わせ"},
		{Name: "大阪環状線", Company: "JR西日本", Status: "大幅な遅延（約20分）"},
		{Name: "阪急神戸線", Company: "阪急電鉄", Status: "平常運転"},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTMLの組み立て（スマホで見やすいカード型のデザイン）
	fmt.Fprintln(w, `
	<!DOCTYPE html>
	<html lang="ja">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>運行情報ダッシュボード</title>
		<!-- 1秒で画面がお洒落になる魔法のCSS -->
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css">
		<style>
			body { padding-top: 20px; background-color: #1a1c1e; }
			.card { margin-bottom: 15px; border-radius: 12px; padding: 15px; background: #2d3135; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
			.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.8rem; font-weight: bold; margin-left: 10px; }
			.badge-ok { background-color: #2e7d32; color: #fff; }
			.badge-ng { background-color: #c62828; color: #fff; }
			.badge-warn { background-color: #f57c00; color: #fff; }
			.company-name { font-size: 0.8rem; color: #a0a5ab; display: block; margin-bottom: 20px; }
		</style>
	</head>
	<body>
		<main class="container" style="max-width: 600px;">
			<h2 style="text-align: center; margin-bottom: 30px;">🚋 リアルタイム運行情報</h2>
	`)

	// 路線ごとのカードを出力
	for _, info := range infoList {
		// ステータスによってバッジの色を変える
		badgeClass := "badge-ok"
		if info.Status == "運転見合わせ" {
			badgeClass = "badge-ng"
		} else if info.Status != "平常運転" {
			badgeClass = "badge-warn"
		}

		fmt.Fprintf(w, `
			<div class="card">
				<span class="company-name">%s</span>
				<div style="display: flex; justify-content: space-between; align-items: center;">
					<strong style="font-size: 1.2rem;">%s</strong>
					<span class="badge %s">%s</span>
				</div>
			</div>
		`, info.Company, info.Name, badgeClass, info.Status)
	}

	fmt.Fprintln(w, `
		</main>
	</body>
	</html>
	`)
}

func main() {
	http.HandleFunc("/delay", delayHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("サーバーを起動しました。ポート: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}