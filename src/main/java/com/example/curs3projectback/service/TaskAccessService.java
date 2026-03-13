package com.example.curs3projectback.service;

import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.model.Task;
import com.example.curs3projectback.model.Tenant;
import com.example.curs3projectback.model.User;
import com.example.curs3projectback.repository.TaskRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class TaskAccessService {

    private final TaskRepository taskRepository;

    public Task getAssignedTask(Long taskId, User assignee, Tenant tenant) {
        return taskRepository.findById(taskId)
                .filter(task -> task.getTenant().getId().equals(tenant.getId()))
                .filter(task -> task.getAssignee() != null && task.getAssignee().getId().equals(assignee.getId()))
                .orElseThrow(() -> new ResourceNotFoundException("Task not found"));
    }
}
