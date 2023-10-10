# クイズアプリ

## 概要
CLI上で完結するインタラクティブな簡易クイズアプリ。  
gRPCの双方向ストリーミングで出題と回答のやり取りを実現している。  
中学生レベルの日本史クイズを50問ほど用意した。  

## 使い方

### クイズの始め方
```
% go run cmd/server/main.go // gRPCサーバーの起動
% go run cmd/client/main.go // クライアント側の起動→クイズスタート

```

### クイズの流れ
1. 「何問出題しますか？」と聞かれるので、自分が回答したい問題数を送信
2. 1問ずつ出題され、その度に回答を入力していく
3. 全ての問題に回答し終えたら、最終結果が送信される

https://github.com/wada1355/quiz-app_grpc/assets/39019484/4c5c91b5-c36f-4b54-8a84-dfbfb35d9301

## TODO
- 入力値バリデーション
- 丁寧なエラーハンドリング
- テスト（単体テスト、シナリオテスト）
- デバッグトレース
