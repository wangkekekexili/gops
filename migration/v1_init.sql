CREATE TABLE IF NOT EXISTS game
(
  `id`        INT NOT NULL AUTO_INCREMENT,
  `name`      VARCHAR(255) NOT NULL,
  `condition` VARCHAR(20) NOT NULL,
  `source`    VARCHAR(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`name`, `condition`, `source`)
) ENGINE=INNODB CHARACTER SET=ascii;

CREATE TABLE IF NOT EXISTS price
(
  `id`        INT NOT NULL AUTO_INCREMENT,
  `game_id`   INT NOT NULL,
  `value`     DECIMAL(6,2) NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  KEY (`game_id`, `timestamp`),
  FOREIGN KEY (`game_id`) REFERENCES game(`id`)
) ENGINE=INNODB CHARACTER SET=ascii;