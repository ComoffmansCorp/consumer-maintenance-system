package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.task.TaskResponse;
import com.example.curs3projectback.dto.task.TaskStatusUpdateRequest;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.mapper.TaskMapper;
import com.example.curs3projectback.model.Task;
import com.example.curs3projectback.model.enums.TaskStatus;
import com.example.curs3projectback.repository.TaskRepository;
import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
public class TaskService {

    private final TaskRepository taskRepository;
    private final UserRepository userRepository;
    private final TaskMapper taskMapper;

    public List<TaskResponse> findMyTasks(TaskStatus status, Authentication authentication) {
        var user = userRepository.findByUsername(authentication.getName())
                .orElseThrow(() -> new ResourceNotFoundException("Пользователь не найден"));
        List<Task> tasks;
        if (status != null) {
            tasks = taskRepository.findAllByAssigneeAndStatus(user, status);
        } else {
            tasks = taskRepository.findAllByAssignee(user);
        }
        return taskMapper.toResponses(tasks);
    }

    public TaskResponse updateStatus(Long taskId, TaskStatusUpdateRequest request, Authentication authentication) {
        var task = getTaskForUser(taskId, authentication);
        task.setStatus(request.getStatus());
        taskRepository.save(task);
        return taskMapper.toResponse(task);
    }

    private Task getTaskForUser(Long taskId, Authentication authentication) {
        var task = taskRepository.findById(taskId)
                .orElseThrow(() -> new ResourceNotFoundException("Задача не найдена"));
        if (!task.getAssignee().getUsername().equals(authentication.getName())) {
            throw new ResourceNotFoundException("Задача не найдена");
        }
        return task;
    }
}

