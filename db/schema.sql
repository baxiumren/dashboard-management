CREATE TABLE IF NOT EXISTS akun_fb (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL DEFAULT '',
    fb_id TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    password TEXT NOT NULL DEFAULT '',
    password_mail TEXT NOT NULL DEFAULT '',
    recovery_mail TEXT NOT NULL DEFAULT '',
    cookie TEXT NOT NULL DEFAULT '',
    twofa_secret TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'aktif',
    tgl_beli TEXT NOT NULL DEFAULT '',
    harga_beli INTEGER NOT NULL DEFAULT 0,
    seller TEXT NOT NULL DEFAULT '',
    catatan TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fanpage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL DEFAULT '',
    page_id TEXT NOT NULL DEFAULT '',
    akun_fb_id INTEGER NOT NULL,
    bm_id INTEGER,
    status TEXT NOT NULL DEFAULT 'aktif',
    tgl_buat TEXT NOT NULL DEFAULT '',
    catatan TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (akun_fb_id) REFERENCES akun_fb(id),
    FOREIGN KEY (bm_id) REFERENCES bm(id)
);

CREATE TABLE IF NOT EXISTS bm (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL DEFAULT '',
    bm_id TEXT NOT NULL DEFAULT '',
    owner_akun_fb_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'aktif',
    tgl_buat TEXT NOT NULL DEFAULT '',
    catatan TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_akun_fb_id) REFERENCES akun_fb(id)
);

CREATE TABLE IF NOT EXISTS bm_shared_akun (
    bm_id INTEGER NOT NULL,
    akun_fb_id INTEGER NOT NULL,
    role TEXT NOT NULL DEFAULT 'employee',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (bm_id, akun_fb_id),
    FOREIGN KEY (bm_id) REFERENCES bm(id),
    FOREIGN KEY (akun_fb_id) REFERENCES akun_fb(id)
);

CREATE TABLE IF NOT EXISTS akun_iklan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL DEFAULT '',
    ad_account_id TEXT NOT NULL DEFAULT '',
    bm_id INTEGER NOT NULL,
    runner_akun_fb_id INTEGER,
    limit_harian INTEGER NOT NULL DEFAULT 0,
    limit_total INTEGER NOT NULL DEFAULT 0,
    metode_bayar TEXT NOT NULL DEFAULT '',
    mata_uang TEXT NOT NULL DEFAULT 'IDR',
    status TEXT NOT NULL DEFAULT 'aktif',
    tgl_buat TEXT NOT NULL DEFAULT '',
    catatan TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bm_id) REFERENCES bm(id),
    FOREIGN KEY (runner_akun_fb_id) REFERENCES akun_fb(id)
);

CREATE TABLE IF NOT EXISTS pixel (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL DEFAULT '',
    pixel_id TEXT NOT NULL DEFAULT '',
    bm_id INTEGER,
    akun_fb_id INTEGER,
    status TEXT NOT NULL DEFAULT 'aktif',
    tgl_buat TEXT NOT NULL DEFAULT '',
    catatan TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bm_id) REFERENCES bm(id),
    FOREIGN KEY (akun_fb_id) REFERENCES akun_fb(id)
);

CREATE TABLE IF NOT EXISTS riwayat (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entitas TEXT NOT NULL,
    entitas_id INTEGER NOT NULL,
    tipe TEXT NOT NULL,
    tanggal TEXT NOT NULL,
    catatan TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
