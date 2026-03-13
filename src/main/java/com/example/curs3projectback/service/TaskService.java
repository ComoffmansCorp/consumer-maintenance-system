package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.task.TaskResponse;
import com.example.curs3projectback.dto.task.TaskStatusUpdateRequest;
import com.example.curs3projectback.mapper.TaskMapper;
import com.example.curs3projectback.model.enums.TaskStatus;
import com.example.curs3projectback.repository.TaskRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class TaskService {

    private final TaskRepository taskRepository;
    private final TaskMapper taskMapper;
    private final CurrentUserService currentUserService;
    private final CurrentTenantService currentTenantService;
    private final TaskAccessService taskAccessService;

    @Transactional(readOnly = true)
    public List<TaskResponse> findMyTasks(TaskStatus status, Authentication authentication) {
        var user = currentUserService.getRequiredUser(authentication);
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Loading tasks userId={} tenant={} status={}", user.getId(), tenant.getCode(), status);
        var tasks = status == null
                ? taskRepository.findAllByTenant_IdAndAssignee(tenant.getId(), user)
                : taskRepository.findAllByTenant_IdAndAssigneeAndStatus(tenant.getId(), user, status);
        return taskMapper.toResponses(tasks);
    }

    @Transactional
    public TaskResponse updateStatus(Long taskId, TaskStatusUpdateRequest request, Authentication authentication) {
        var user = currentUserService.getRequiredUser(authentication);
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Updating task status taskId={} tenant={} userId={} newStatus={}",
                taskId, tenant.getCode(), user.getId(), request.getStatus());
        var task = taskAccessService.getAssignedTask(taskId, user, tenant);
        task.setStatus(request.getStatus());
        return taskMapper.toResponse(taskRepository.save(task));
    }
}
