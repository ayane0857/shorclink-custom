async function createUser() {
  try {
    const response = await fetch("http://localhost:8080/api/7", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "X-API-Token": "your_api_token",
      },
      body: JSON.stringify({
        url: "https://music.youtube.com/watch?v=TG_ZJLL8YXYYYYYYYYYYYY",
      }),
    });

    // レスポンスの中身をテキストで取得（JSONでない可能性もあるので）
    const text = await response.text();

    if (!response.ok) {
      console.error("=== HTTPエラー詳細 ===");
      console.error("Status:", response.status);
      console.error("Body:", text);
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    // JSONとしてパースできる場合だけ
    const data = JSON.parse(text);
    console.log("=== 成功 ===");
    console.log(data);
  } catch (err) {
    console.error("=== Fetchで例外発生 ===");
    console.error(err);
  }
}

createUser();
