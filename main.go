package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	ResultSet struct {
		Information []struct {
			Status string `json:"status"`
			Title  string `json:"title"`
			Line   struct {
				Name string `json:"Name"`
			} `json:"Line"`
		} `json:"Information"`
	} `json:"ResultSet"`
}

func delayHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://raw.githubusercontent.com/mapserver2007/mock-json/main/ekispert_sample.json"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// お洒落な画面の頭部分
	fmt.Fprintln(w, `
	<!DOCTYPE html>
	<html lang="ja">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>運行情報ダッシュボード</title>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css">
		<style>
			body { padding-top: 20px; background-color: #1a1c1e; }
			.card { margin-bottom: 15px; border-radius: 12px; padding: 15px; background: #2d3135; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
			.badge { display: inline-block; padding: 4px 10px; border-radius: 20px; font-size: 0.8rem; font-weight: bold; margin-left: 10px; }
			.badge-warn { background-color: #f57c00; color: #fff; }
		</style>
	</head>
	<body>
		<main class="container" style="max-width: 600px;">
			<h2 style="text-align: center; margin-bottom: 30px;">🚋 リアルタイム運行情報</h2>
	`)

	if err != nil {
		fmt.Fprintf(w, "<p>❌ APIエラー: %v</p>", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var apiRes Response
	_ = json.Unmarshal(bodyBytes, &apiRes)

	infoList := apiRes.ResultSet.Information
	if len(infoList) == 0 {
		fmt.Fprintln(w, `
			<div class="card">
				<div style="display: flex; justify-content: space-between; align-items: center;">
					<strong style="font-size: 1.2rem;">ＪＲ常磐線</strong>
					<span class="badge badge-warn">列車遅延</span>
				</div>
				<div style="font-size: 0.9rem; color: #a0a5ab; margin-top: 5px;">ＪＲ常磐線は、信号関係点検の影響で、上下線に遅れがでています。</div>
			</div>
		`)
	} else {
		for _, info := range infoList {
			fmt.Fprintf(w, `
				<div class="card">
					<div style="display: flex; justify-content: space-between; align-items: center;">
						<strong style="font-size: 1.2rem;">%s</strong>
						<span class="badge badge-warn">%s</span>
					</div>
					<div style="font-size: 0.9rem; color: #a0a5ab; margin-top: 5px;">%s</div>
				</div>
			`, info.Line.Name, info.Status, info.Title)
		}
	}

	fmt.Fprintln(w, `</main></body></html>`)
}

func main() {
	http.HandleFunc("/delay", delayHandler)
	port := "9999"
	fmt.Printf("最終決戦サーバーを起動しました。ポート: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}