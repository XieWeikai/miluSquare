-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema miluSquarePro
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema miluSquarePro
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `miluSquarePro` DEFAULT CHARACTER SET utf8 ;
USE `miluSquarePro` ;

-- -----------------------------------------------------
-- Table `miluSquarePro`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `miluSquarePro`.`users` (
  `user_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_name` VARCHAR(45) NOT NULL DEFAULT '无名',
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `user_id_UNIQUE` (`user_id` ASC) VISIBLE,
  UNIQUE INDEX `email_UNIQUE` (`email` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `miluSquarePro`.`communities`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `miluSquarePro`.`communities` (
  `community_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL DEFAULT '麋鹿社区',
  `description` VARCHAR(255) NOT NULL DEFAULT '暂无描述',
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `belong_to` VARCHAR(45) NOT NULL DEFAULT '无所属',
  `user_id` INT UNSIGNED NOT NULL,
  PRIMARY KEY (`community_id`),
  UNIQUE INDEX `community_id_UNIQUE` (`community_id` ASC) VISIBLE,
  INDEX `fk_communities_users1_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_communities_users1`
    FOREIGN KEY (`user_id`)
    REFERENCES `miluSquarePro`.`users` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE CASCADE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `miluSquarePro`.`posts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `miluSquarePro`.`posts` (
  `post_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL DEFAULT '无标题',
  `topic` VARCHAR(45) NOT NULL DEFAULT '无主题',
  `content` MEDIUMTEXT NOT NULL,
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `community_id` INT UNSIGNED NOT NULL,
  `user_id` INT UNSIGNED NOT NULL,
  PRIMARY KEY (`post_id`),
  UNIQUE INDEX `post_id_UNIQUE` (`post_id` ASC) VISIBLE,
  INDEX `fk_posts_communities_idx` (`community_id` ASC) VISIBLE,
  INDEX `fk_posts_users1_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_posts_communities`
    FOREIGN KEY (`community_id`)
    REFERENCES `miluSquarePro`.`communities` (`community_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_posts_users1`
    FOREIGN KEY (`user_id`)
    REFERENCES `miluSquarePro`.`users` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE CASCADE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `miluSquarePro`.`comments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `miluSquarePro`.`comments` (
  `comment_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `content` TEXT NULL,
  `post_id` INT UNSIGNED NOT NULL,
  `user_id` INT UNSIGNED NOT NULL,
  PRIMARY KEY (`comment_id`),
  UNIQUE INDEX `comment_id_UNIQUE` (`comment_id` ASC) VISIBLE,
  INDEX `fk_comments_posts1_idx` (`post_id` ASC) VISIBLE,
  INDEX `fk_comments_users1_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_comments_posts1`
    FOREIGN KEY (`post_id`)
    REFERENCES `miluSquarePro`.`posts` (`post_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_comments_users1`
    FOREIGN KEY (`user_id`)
    REFERENCES `miluSquarePro`.`users` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE CASCADE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `miluSquarePro`.`imgs`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `miluSquarePro`.`imgs` (
  `img_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `path` VARCHAR(255) NOT NULL,
  `post_id` INT UNSIGNED NOT NULL,
  PRIMARY KEY (`img_id`),
  UNIQUE INDEX `img_id_UNIQUE` (`img_id` ASC) VISIBLE,
  INDEX `fk_imgs_posts1_idx` (`post_id` ASC) VISIBLE,
  CONSTRAINT `fk_imgs_posts1`
    FOREIGN KEY (`post_id`)
    REFERENCES `miluSquarePro`.`posts` (`post_id`)
    ON DELETE NO ACTION
    ON UPDATE CASCADE)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
