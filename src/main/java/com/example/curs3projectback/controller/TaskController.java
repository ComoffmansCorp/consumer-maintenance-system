package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.task.TaskResponse;
import com.example.curs3projectback.dto.task.TaskStatusUpdateRequest;
import com.example.curs3projectback.model.enums.TaskStatus;
import com.example.curs3projectback.service.TaskService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/api/tasks")
@RequiredArgsConstructor
public class TaskController {

    private final TaskService taskService;

    @GetMapping
    public ResponseEntity<List<TaskResponse>> myTasks(@RequestParam(value = "status", required = false) TaskStatus status,
                                                      Authentication authentication) {
        log.info("HTTP GET /api/tasks userId={} status={}", authentication.getName(), status);
        return ResponseEntity.ok(taskService.findMyTasks(status, authentication));
    }

    @PutMapping("/{taskId}/status")
    public ResponseEntity<TaskResponse> updateStatus(@PathVariable Long taskId,
                                                     @Valid @RequestBody TaskStatusUpdateRequest request,
                                                     Authentication authentication) {
        log.info("HTTP PUT /api/tasks/{}/status userId={} newStatus={}", taskId, authentication.getName(), request.getStatus());
        return ResponseEntity.ok(taskService.updateStatus(taskId, request, authentication));
    }
}

