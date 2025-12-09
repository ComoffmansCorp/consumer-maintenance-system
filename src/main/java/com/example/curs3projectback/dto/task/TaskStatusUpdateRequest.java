package com.example.curs3projectback.dto.task;

import com.example.curs3projectback.model.enums.TaskStatus;
import lombok.Data;

@Data
public class TaskStatusUpdateRequest {
    private TaskStatus status;
}

