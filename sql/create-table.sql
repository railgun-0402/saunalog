use sauna-db;
CREATE TABLE IF NOT EXISTS `users` (
    `id`         INT PRIMARY KEY AUTO_INCREMENT,
    `name`       varchar(20) NOT NULL,
    `email`      varchar(255) NOT NULL UNIQUE,
    `gender`     varchar(10) NOT NULL,
    `age`        INT NOT NULL,
    `password`   varchar(255) NOT NULL,
    `prefecture` varchar(20) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
