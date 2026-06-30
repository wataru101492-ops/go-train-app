package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type TrainInfo struct {
	Name    string
	Company string
	Status  string
}

func delayHandler(w http.ResponseWriter, r *http.Request) {

	infoList := []TrainInfo{
		{Name: "JR京都線", Company: "JR西日本", Status: "平常運転"},
		{Name: "御堂筋線", Company: "大阪メトロ", Status: "運転見合わせ"},
		{Name: "大阪環状線", Company: "JR西日本", Status: "約20分遅延"},
		{Name: "阪急神戸線", Company: "阪急電鉄", Status: "平常運転"},
	}

	ok := 0
	warn := 0
	stop := 0

	for _, i := range infoList {
		switch {
		case i.Status == "平常運転":
			ok++
		case i.Status == "運転見合わせ":
			stop++
		default:
			warn++
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">

<title>運行情報</title>

<style>

*{
    box-sizing:border-box;
}

body{
    margin:0;
    font-family:
        "Hiragino Sans",
        "Noto Sans JP",
        sans-serif;

    background:
        linear-gradient(180deg,#0f172a,#1e293b);

    color:white;
}

.header{
    padding:30px 20px;
    text-align:center;
}

.header h1{
    margin:0;
    font-size:1.8rem;
}

.updated{
    margin-top:10px;
    color:#94a3b8;
    font-size:0.9rem;
}

.summary{

    display:flex;
    justify-content:space-around;

    margin:20px;

    padding:18px;

    border-radius:18px;

    background:#1e293b;

    box-shadow:
      0 8px 25px rgba(0,0,0,.35);

}

.summary-item{
    text-align:center;
}

.summary-num{
    font-size:1.5rem;
    font-weight:bold;
}

.container{
    padding:0 15px 30px;
}

.card{

    background:#1e293b;

    border-radius:18px;

    padding:18px;

    margin-bottom:15px;

    display:flex;

    justify-content:space-between;

    align-items:center;

    box-shadow:
      0 8px 20px rgba(0,0,0,.35);

}

.ok{
    border-left:6px solid #22c55e;
}

.warn{
    border-left:6px solid #f59e0b;
}

.stop{
    border-left:6px solid #ef4444;
}

.company{
    color:#94a3b8;
    font-size:0.8rem;
}

.line{

    margin-top:5px;

    font-size:1.15rem;

    font-weight:bold;

}

.badge{

    padding:8px 12px;

    border-radius:999px;

    font-size:.82rem;

    font-weight:bold;

    white-space:nowrap;

}

.badge-ok{
    background:#16a34a;
}

.badge-warn{
    background:#f59e0b;
}

.badge-stop{
    background:#dc2626;
}

</style>

</head>

<body>

<div class="header">

<h1>🚋 運行情報</h1>

<div class="updated">更新 `)

	fmt.Fprintf(w, "%s", time.Now().Format("15:04"))

	fmt.Fprint(w, `</div>

</div>

<div class="summary">

<div class="summary-item">
<div class="summary-num">🟢 `)

	fmt.Fprintf(w, "%d", ok)

	fmt.Fprint(w, `</div>
<div>平常</div>
</div>

<div class="summary-item">
<div class="summary-num">🟠 `)

	fmt.Fprintf(w, "%d", warn)

	fmt.Fprint(w, `</div>
<div>遅延</div>
</div>

<div class="summary-item">
<div class="summary-num">🔴 `)

	fmt.Fprintf(w, "%d", stop)

	fmt.Fprint(w, `</div>
<div>停止</div>
</div>

</div>

<div class="container">
`)

	for _, info := range infoList {

		card := "ok"
		badge := "badge-ok"
		icon := "✅"

		switch {
		case info.Status == "平常運転":
			card = "ok"
			badge = "badge-ok"
			icon = "✅"

		case info.Status == "運転見合わせ":
			card = "stop"
			badge = "badge-stop"
			icon = "⛔"

		default:
			card = "warn"
			badge = "badge-warn"
			icon = "⚠️"
		}

		fmt.Fprintf(w, `
<div class="card %s">

<div>
<div class="company">%s</div>
<div class="line">🚆 %s</div>
</div>

<div class="badge %s">
%s %s
</div>

</div>
`,
			card,
			info.Company,
			info.Name,
			badge,
			icon,
			info.Status,
		)
	}

	fmt.Fprint(w, `

</div>

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

	fmt.Println("http://localhost:" + port + "/delay")

	http.ListenAndServe(":"+port, nil)
}