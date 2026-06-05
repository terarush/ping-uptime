#!/bin/bash

TARGET_DIR="/usr/local/bin"

echo "============================================="
echo " Installing ping-uptime (Auto-detecting version)"
echo "============================================="

# 1. Mengambil versi tag terbaru secara otomatis dari GitHub API
echo "🔍 Mendeteksi versi rilis terbaru dari GitHub..."
VERSION=$(curl -s https://api.github.com/repos/terarush/ping-uptime/releases/latest | grep '"tag_name":' | cut -d '"' -f 4)

if [ -z "$VERSION" ]; then
    echo "⚠️  Gagal mendeteksi versi secara otomatis. Menggunakan fallback v0.0.3..."
    VERSION="v0.0.3"
else
    echo "✅ Versi terbaru terdeteksi: $VERSION"
fi

# Link download dinamis menggunakan versi yang terdeteksi
DOWNLOAD_URL="https://github.com/terarush/ping-uptime/releases/download/${VERSION}/ping-uptime"

# 2. Download binary ke folder tmp
echo "📥 Downloading binary dari GitHub..."
curl -L "$DOWNLOAD_URL" -o /tmp/ping-uptime

if [ $? -ne 0 ]; then
    echo "❌ Gagal mengunduh file. Periksa koneksi internet Anda."
    exit 1
fi

# 3. Buat file menjadi executable
chmod +x /tmp/ping-uptime

# 4. Pindahkan ke /usr/local/bin
echo "⚙️  Memindahkan binary ke $TARGET_DIR..."
sudo mv /tmp/ping-uptime "$TARGET_DIR/ping-uptime"

if [ $? -ne 0 ]; then
    echo "❌ Gagal memindahkan file. Pastikan Anda memiliki akses sudo."
    exit 1
fi

# 5. Deteksi Shell dan update PATH di .bashrc atau .zshrc
DETECTED_SHELL=$(basename "$SHELL")
RC_FILE=""

if [ "$DETECTED_SHELL" = "zsh" ]; then
    RC_FILE="$HOME/.zshrc"
elif [ "$DETECTED_SHELL" = "bash" ]; then
    RC_FILE="$HOME/.bashrc"
fi

if [ -n "$RC_FILE" ] && [ -f "$RC_FILE" ]; then
    # Cek apakah /usr/local/bin sudah ada di PATH rc file
    if ! grep -q "$TARGET_DIR" "$RC_FILE"; then
        echo "📝 Menambahkan $TARGET_DIR ke PATH di $RC_FILE..."
        echo "export PATH=\"$TARGET_DIR:\$PATH\"" >> "$RC_FILE"
        echo "✅ Konfigurasi shell berhasil di-update."
        echo "👉 Silakan jalankan: source $RC_FILE untuk me-refresh terminal Anda."
    else
        echo "ℹ️  PATH $TARGET_DIR sudah terdaftar di $RC_FILE."
    fi
else
    echo "⚠️  Shell tidak dikenal atau file konfigurasi tidak ditemukan. Pastikan $TARGET_DIR ada di PATH Anda secara manual."
fi

echo "============================================="
echo "🎉 Instalasi Selesai! Coba jalankan: ping-uptime -h"
echo "============================================="
