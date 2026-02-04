/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package server

import "fmt"

// generateHTML generates the HTML page for file download
func GenerateHTML(fileName string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LAN Share - %s</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 25%%, #f093fb 50%%, #4facfe 75%%, #667eea 100%%);
            background-size: 400%% 400%%;
            animation: gradientShift 15s ease infinite;
            padding: 20px;
        }

        @keyframes gradientShift {
            0%% { background-position: 0%% 50%%; }
            50%% { background-position: 100%% 50%%; }
            100%% { background-position: 0%% 50%%; }
        }

        .container {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 24px;
            padding: 48px;
            max-width: 500px;
            width: 100%%;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            text-align: center;
            animation: fadeIn 0.6s ease-out;
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .icon {
            width: 80px;
            height: 80px;
            margin: 0 auto 24px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            border-radius: 20px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 40px;
            box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
        }

        h1 {
            color: #2d3748;
            font-size: 28px;
            font-weight: 700;
            margin-bottom: 12px;
        }

        .subtitle {
            color: #718096;
            font-size: 14px;
            margin-bottom: 32px;
            text-transform: uppercase;
            letter-spacing: 1px;
            font-weight: 600;
        }

        .file-name {
            background: linear-gradient(135deg, #f6f8fb 0%%, #e9ecef 100%%);
            padding: 20px;
            border-radius: 12px;
            margin-bottom: 32px;
            word-break: break-all;
            border: 2px solid #e2e8f0;
        }

        .file-name-text {
            color: #2d3748;
            font-size: 16px;
            font-weight: 600;
            font-family: 'Courier New', monospace;
        }

        .download-btn {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            border: none;
            padding: 16px 48px;
            font-size: 18px;
            font-weight: 600;
            border-radius: 12px;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
            text-decoration: none;
            display: inline-block;
        }

        .download-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 12px 28px rgba(102, 126, 234, 0.4);
        }

        .download-btn:active {
            transform: translateY(0);
        }

        .footer {
            margin-top: 32px;
            color: #a0aec0;
            font-size: 13px;
        }

        @media (max-width: 600px) {
            .container {
                padding: 32px 24px;
            }

            h1 {
                font-size: 24px;
            }

            .download-btn {
                width: 100%%;
                padding: 14px 32px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">📁</div>
        <h1>File Ready to Download</h1>
        <p class="subtitle">LAN Share</p>
        
        <div class="file-name">
            <div class="file-name-text">%s</div>
        </div>
        
        <a href="/download" class="download-btn">⬇️ Download File</a>
        
        <div class="footer">
            Click the button above to download the file
        </div>
    </div>
</body>
</html>`, fileName, fileName)
}
