cd ssl
IF EXIST ssl_generator.sh (
    ren ssl_generator.sh ssl_generator.bat
    ssl_generator.bat
)
cd ..

go run .