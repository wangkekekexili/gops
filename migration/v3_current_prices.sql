CREATE TABLE IF NOT EXISTS current_prices (
    `id` INT NOT NULL AUTO_INCREMENT,
    `game_id` INT NOT NULL,
    `price_id` INT NOT NULL,
    `value` DECIMAL(6, 2) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`game_id`),
    FOREIGN KEY (`game_id`) REFERENCES games(`id`),
    FOREIGN KEY (`price_id`) REFERENCES prices(`id`)
) ENGINE=INNODB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci;
