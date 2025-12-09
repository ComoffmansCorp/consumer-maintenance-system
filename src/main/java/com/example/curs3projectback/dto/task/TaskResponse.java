package com.example.curs3projectback.dto.task;

import com.example.curs3projectback.model.enums.TaskStatus;
import com.example.curs3projectback.model.enums.TaskType;
import lombok.Builder;
import lombok.Value;

import java.time.LocalDate;

@Value
@Builder
public class TaskResponse {
    Long id;
    TaskType type;
    String address;
    String consumerName;
    TaskStatus status;
    LocalDate dueDate;
}

