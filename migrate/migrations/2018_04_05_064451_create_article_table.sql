-- Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `article` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `age` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `article`;
