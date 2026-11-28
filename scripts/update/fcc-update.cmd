@echo off
rem Do not download the upstream installer or infer a checkout to modify.
>&2 echo Automatic update is not available for the MoLik-Labs/FCC source checkout.
>&2 echo The upstream installer would update a different project; nothing was changed.
>&2 echo Stop FCC, inspect local changes, then run in your FCC checkout:
>&2 echo   git pull --ff-only
>&2 echo   uv sync --locked
>&2 echo If using local voice, include --extra voice_local when syncing.
exit /b 2
