INSERT INTO `users` (`name`, `password`) VALUES ("user", unhex("c7d45c33b66657620d2327fa216d91c6e1970d92e0731898c552885e13bca53c"));
INSERT INTO `users` (`name`, `password`) VALUES ("other", unhex("c4f6e0d91770dba45243f6b781c568196c1832725414a0d3f8271f0e492e6085"));

INSERT INTO `tasks` (`title`, `description`) VALUES ("sample-task-01", "task 01");
INSERT INTO `tasks` (`title`) VALUES ("sample-task-02");
INSERT INTO `tasks` (`title`, `description`)
VALUES ("sample-task-03", "task 03 long long long long long long long long description");
INSERT INTO `tasks` (`title`) VALUES ("sample-task-04");
INSERT INTO `tasks` (`title`, `is_done`) VALUES ("sample-task-05", TRUE);

INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (1, 1);
INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (1, 2);

INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (1, 3);
INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (2, 4);
INSERT INTO `ownership` (`user_id`, `task_id`) VALUES (1, 5);
