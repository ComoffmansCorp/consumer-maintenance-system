package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.Task;
import com.example.curs3projectback.model.User;
import com.example.curs3projectback.model.enums.TaskStatus;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface TaskRepository extends JpaRepository<Task, Long> {
    List<Task> findAllByAssigneeAndStatus(User assignee, TaskStatus status);
    List<Task> findAllByAssignee(User assignee);
}

