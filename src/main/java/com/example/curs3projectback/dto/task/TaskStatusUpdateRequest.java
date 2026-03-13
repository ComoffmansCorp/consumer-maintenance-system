package com.example.curs3projectback.dto.task;

import com.example.curs3projectback.model.enums.TaskStatus;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

@Data
public class TaskStatusUpdateRequest {
    @NotNull(message = "status is required")
    private TaskStatus status;
}

