go build -o judge/judge github.com/ggaaooppeenngg/OJ/judge/main.go
judge/judge &
echo "begin judge loop"
revel run github.com/ggaaooppeenngg/OJ
