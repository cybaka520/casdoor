package ee

import (
	"fmt"
	"html"
	"strings"
	"time"
)

const (
	alertBrandName = "AMLL Hub"
	alertBrandSite = "https://beta.amllhub.bikonoo.com"
	alertBrandLogo = "https://beta.amllhub.bikonoo.com/logo.png"
)

// signinAlertRow 登录提醒里的一行明细
type signinAlertRow struct {
	Label string
	Value string
}

// buildSigninRiskEmail 生成与 AMLL Hub 站内邮件一致的 HTML 登录提醒。
// 正文里的所有动态值都已做 HTML 转义，避免用户可控的 UA / 设备名注入标签。
func buildSigninRiskEmail(displayName, userName string, loginTime time.Time, rows []signinAlertRow) string {
	greeting := displayName
	if greeting == "" {
		greeting = userName
	}

	var details strings.Builder
	for _, row := range rows {
		details.WriteString(fmt.Sprintf(`
          <tr>
            <td style="padding:7px 0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:13px;line-height:20px;color:#8a8a90;width:64px;vertical-align:top;">%s</td>
            <td style="padding:7px 0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:14px;line-height:22px;color:#1d1d1f;word-break:break-all;">%s</td>
          </tr>`,
			html.EscapeString(row.Label), html.EscapeString(row.Value)))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="color-scheme" content="light dark">
<title>新登录提醒</title>
</head>
<body style="margin:0;padding:0;background:#fbfbfd;">
<div style="display:none;font-size:1px;line-height:1px;max-height:0;max-width:0;opacity:0;overflow:hidden;">
  你的 AMLL Hub 账号有新的登录活动。
</div>
<table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background:#fbfbfd;">
<tr>
<td align="center" style="padding:28px 12px 36px 12px;">
  <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="600" style="width:600px;max-width:600px;">

    <tr>
      <td style="padding:0 4px 18px 4px;">
        <a href="%[1]s" style="text-decoration:none;">
          <img src="%[2]s" width="32" height="32" alt="AMLL Hub" style="display:inline-block;border:0;border-radius:8px;vertical-align:middle;">
          <span style="display:inline-block;vertical-align:middle;margin-left:8px;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:16px;font-weight:700;letter-spacing:.3px;color:#1d1d1f;">AMLL Hub</span>
        </a>
      </td>
    </tr>

    <tr>
      <td style="padding:0 4px 16px 4px;">
        <div style="height:1px;line-height:1px;font-size:0;background:#ededf1;">&nbsp;</div>
      </td>
    </tr>

    <tr>
      <td style="background:#ffffff;border:1px solid #ededf1;border-radius:16px;padding:34px 40px 32px 40px;">
        <table role="presentation" border="0" cellpadding="0" cellspacing="0">
          <tr>
            <td style="background:#fbe7e9;border-radius:999px;padding:7px 15px;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:13px;font-weight:600;color:#e0303f;line-height:20px;">
              安全提醒
            </td>
          </tr>
        </table>

        <h1 style="margin:16px 0 4px 0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:22px;line-height:31px;font-weight:600;color:#1d1d1f;">
          检测到新的登录
        </h1>
        <p style="margin:0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:13px;line-height:22px;color:#8a8a90;">
          %[3]s　·　%[4]s
        </p>

        <p style="margin:22px 0 10px 0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:15px;line-height:26px;color:#535359;">
          %[5]s，你的账号刚刚完成了一次登录：
        </p>

        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%" style="background:#fbfbfd;border:1px solid #ededf1;border-radius:12px;">
          <tr>
            <td style="padding:10px 18px;">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%">
                %[6]s
              </table>
            </td>
          </tr>
        </table>

        <p style="margin:20px 0 0 0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:14px;line-height:24px;color:#535359;">
          如果这不是你本人的操作，请立即修改密码并退出所有设备；如有疑问，请联系管理员。
        </p>

        <table role="presentation" border="0" cellpadding="0" cellspacing="0" style="margin-top:22px;">
          <tr>
            <td align="center" bgcolor="#e0303f" style="background-color:#e0303f;border-radius:8px;">
              <a href="%[1]s/reset-password" style="display:inline-block;padding:13px 30px;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:15px;font-weight:600;color:#ffffff;text-decoration:none;border-radius:8px;">修改密码</a>
            </td>
          </tr>
        </table>
      </td>
    </tr>

    <tr>
      <td align="center" style="padding:24px 20px 0 20px;">
        <p style="margin:0 0 6px 0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:12px;line-height:20px;color:#8a8a90;">
          本邮件由 AMLL Hub 系统自动发送，请勿直接回复。
        </p>
        <p style="margin:0;font-family:Inter,'PingFang SC','Noto Sans SC','Microsoft YaHei',system-ui,sans-serif;font-size:12px;line-height:20px;color:#8a8a90;">
          <a href="%[1]s" style="color:#e0303f;text-decoration:none;">%[1]s</a>
          　·　© %[7]d %[8]s
        </p>
      </td>
    </tr>

  </table>
</td>
</tr>
</table>
</body>
</html>`,
		alertBrandSite,
		alertBrandLogo,
		alertBrandName,
		loginTime.Format("2006-01-02 15:04:05"),
		html.EscapeString(greeting),
		details.String(),
		loginTime.Year(),
		alertBrandName,
	)
}
