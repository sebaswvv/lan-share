/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package server

// generateUploadHTML generates the HTML page for file uploads
func GenerateUploadHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LAN Share - Upload</title>
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
            background: linear-gradient(135deg, #667eea 0%, #764ba2 25%, #f093fb 50%, #4facfe 75%, #667eea 100%);
            background-size: 400% 400%;
            animation: gradientShift 15s ease infinite;
            padding: 20px;
        }

        @keyframes gradientShift {
            0% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
            100% { background-position: 0% 50%; }
        }

        .container {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 24px;
            padding: 48px;
            max-width: 500px;
            width: 100%;
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
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

        .upload-area {
            border: 3px dashed #cbd5e0;
            border-radius: 12px;
            padding: 48px 24px;
            margin-bottom: 24px;
            cursor: pointer;
            transition: all 0.3s ease;
            background: #f7fafc;
        }

        .upload-area:hover {
            border-color: #667eea;
            background: #edf2f7;
        }

        .upload-area.drag-over {
            border-color: #667eea;
            background: #e6f2ff;
        }

        .upload-icon {
            font-size: 48px;
            margin-bottom: 16px;
        }

        .upload-text {
            color: #4a5568;
            font-size: 16px;
            margin-bottom: 8px;
        }

        .upload-hint {
            color: #a0aec0;
            font-size: 14px;
        }

        input[type="file"] {
            display: none;
        }

        .selected-file {
            background: linear-gradient(135deg, #f6f8fb 0%, #e9ecef 100%);
            padding: 16px;
            border-radius: 12px;
            margin-bottom: 24px;
            display: none;
        }

        .selected-file.show {
            display: block;
        }

        .file-info {
            color: #2d3748;
            font-weight: 600;
        }

        .upload-btn {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 16px 48px;
            font-size: 18px;
            font-weight: 600;
            border-radius: 12px;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 8px 20px rgba(102, 126, 234, 0.3);
            width: 100%;
            display: none;
        }

        .upload-btn.show {
            display: block;
        }

        .upload-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 12px 28px rgba(102, 126, 234, 0.4);
        }

        .upload-btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            transform: none;
        }

        .progress {
            display: none;
            margin-top: 24px;
        }

        .progress.show {
            display: block;
        }

        .progress-bar {
            width: 100%;
            height: 8px;
            background: #e2e8f0;
            border-radius: 4px;
            overflow: hidden;
            margin-bottom: 8px;
        }

        .progress-fill {
            height: 100%;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            width: 0%;
            transition: width 0.3s ease;
        }

        .progress-text {
            color: #718096;
            font-size: 14px;
        }

        @media (max-width: 600px) {
            .container {
                padding: 32px 24px;
            }

            h1 {
                font-size: 24px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">📤</div>
        <h1>Upload File</h1>
        <p class="subtitle">LAN Share</p>
        
        <form id="uploadForm">
            <div class="upload-area" id="uploadArea">
                <div class="upload-icon">📁</div>
                <div class="upload-text">Click to select or drag & drop</div>
                <div class="upload-hint">Any file type supported</div>
            </div>
            
            <input type="file" id="fileInput" name="file">
            
            <div class="selected-file" id="selectedFile">
                <div class="file-info" id="fileInfo"></div>
            </div>
            
            <button type="submit" class="upload-btn" id="uploadBtn">
                ⬆️ Upload File
            </button>
        </form>
        
        <div class="progress" id="progress">
            <div class="progress-bar">
                <div class="progress-fill" id="progressFill"></div>
            </div>
            <div class="progress-text" id="progressText">Uploading...</div>
        </div>
    </div>

    <script>
        const uploadArea = document.getElementById('uploadArea');
        const fileInput = document.getElementById('fileInput');
        const selectedFile = document.getElementById('selectedFile');
        const fileInfo = document.getElementById('fileInfo');
        const uploadBtn = document.getElementById('uploadBtn');
        const uploadForm = document.getElementById('uploadForm');
        const progress = document.getElementById('progress');
        const progressFill = document.getElementById('progressFill');
        const progressText = document.getElementById('progressText');

        uploadArea.addEventListener('click', () => fileInput.click());

        fileInput.addEventListener('change', (e) => {
            if (e.target.files.length > 0) {
                const file = e.target.files[0];
                showSelectedFile(file);
            }
        });

        uploadArea.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadArea.classList.add('drag-over');
        });

        uploadArea.addEventListener('dragleave', () => {
            uploadArea.classList.remove('drag-over');
        });

        uploadArea.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadArea.classList.remove('drag-over');
            
            if (e.dataTransfer.files.length > 0) {
                fileInput.files = e.dataTransfer.files;
                showSelectedFile(e.dataTransfer.files[0]);
            }
        });

        function showSelectedFile(file) {
            const sizeMB = (file.size / (1024 * 1024)).toFixed(2);
            fileInfo.textContent = file.name + ' (' + sizeMB + ' MB)';
            selectedFile.classList.add('show');
            uploadBtn.classList.add('show');
        }

        uploadForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            if (fileInput.files.length === 0) return;

            const formData = new FormData();
            formData.append('file', fileInput.files[0]);

            uploadBtn.disabled = true;
            progress.classList.add('show');

            try {
                const xhr = new XMLHttpRequest();
                
                xhr.upload.addEventListener('progress', (e) => {
                    if (e.lengthComputable) {
                        const percent = (e.loaded / e.total) * 100;
                        progressFill.style.width = percent + '%';
                        progressText.textContent = 'Uploading... ' + Math.round(percent) + '%';
                    }
                });

                xhr.addEventListener('load', () => {
                    if (xhr.status === 200) {
                        document.body.innerHTML = xhr.responseText;
                    } else {
                        alert('Upload failed!');
                        uploadBtn.disabled = false;
                        progress.classList.remove('show');
                    }
                });

                xhr.addEventListener('error', () => {
                    alert('Upload error!');
                    uploadBtn.disabled = false;
                    progress.classList.remove('show');
                });

                xhr.open('POST', '/upload');
                xhr.send(formData);
            } catch (error) {
                alert('Upload failed!');
                uploadBtn.disabled = false;
                progress.classList.remove('show');
            }
        });
    </script>
</body>
</html>`
}
