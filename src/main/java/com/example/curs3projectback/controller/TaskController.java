package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.task.TaskResponse;
import com.example.curs3projectback.dto.task.TaskStatusUpdateRequest;
import com.example.curs3projectback.model.enums.TaskStatus;
import com.example.curs3projectback.service.TaskService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/tasks")
@RequiredArgsConstructor
public class TaskController {

    private final TaskService taskService;

    @GetMapping
    public ResponseEntity<List<TaskResponse>> myTasks(@RequestParam(value = "status", required = false) TaskStatus status,
                                                      Authentication authentication) {
        return ResponseEntity.ok(taskService.findMyTasks(status, authentication));
    }

    @PutMapping("/{taskId}/status")
    public ResponseEntity<TaskResponse> updateStatus(@PathVariable Long taskId,
                                                     @RequestBody TaskStatusUpdateRequest request,
                                                     Authentication authentication) {
        return ResponseEntity.ok(taskService.updateStatus(taskId, request, authentication));
    }
}

