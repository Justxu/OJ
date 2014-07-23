go build -o judge/judge github.com/ggaaooppeenngg/OJ/judge
judge/judge &
echo "begin judge loop"
revel run github.com/ggaaooppeenngg/OJ
