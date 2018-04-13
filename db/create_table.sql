DROP TABLE IF EXISTS `oauth_user`;

CREATE TABLE `oauth_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(255) unsigned NOT NULL,
  `acccess_token` varchar(40) DEFAULT NULL,
  `refresh_token` varchar(40) DEFAULT NULL,
  `expiry_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int(255) unsigned NOT NULL AUTO_INCREMENT,
  `student_id` varchar(8) NOT NULL DEFAULT '',
  `password` varchar(16) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  `sex` int(2) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
