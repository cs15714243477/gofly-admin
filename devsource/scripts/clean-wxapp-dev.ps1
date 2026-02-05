$ErrorActionPreference = "Stop"

Set-Location -Path (Resolve-Path (Join-Path $PSScriptRoot "..\\.."))

$target = "web\\wxapp\\unpackage\\dist\\dev"
if (Test-Path $target) {
  $maxRetry = 5
  for ($i = 1; $i -le $maxRetry; $i++) {
    try {
      Remove-Item -Recurse -Force $target -ErrorAction Stop
      Write-Host "已删除：$target"
      exit 0
    } catch {
      if ($i -lt $maxRetry) {
        Start-Sleep -Seconds 1
        continue
      }
      Write-Host "删除失败：$target"
      Write-Host "原因：可能有进程正在占用文件（例如微信开发者工具/真机调试/IDE 预览）。"
      Write-Host "处理方式：请先关闭微信开发者工具、停止真机调试，然后重试执行本脚本。"
      Write-Host ("错误详情：" + $_.Exception.Message)
      exit 1
    }
  }
} else {
  Write-Host "未找到：$target（无需删除）"
}

