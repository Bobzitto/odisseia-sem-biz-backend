@echo off
echo Restoring db from create_tables.sql
set DB_PASS=jb52vWqkwxxEiGsC55nvE6aLgrbGL63Z
set DB_USER="odisseiasembiz_user"
set DB_NAME="odisseiasembiz"
set DUMP_FILE="create_tables.sql"
set HOSTNAME="dpg-cs3uicogph6c73c8cnd0-a.oregon-postgres.render.com"

:: Restore the database
set PGPASSWORD=%DB_PASS%
psql -h %HOSTNAME% -U %DB_USER% -d %DB_NAME% -f %DUMP_FILE%
pause