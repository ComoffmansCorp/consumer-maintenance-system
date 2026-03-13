package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.act.MeterInfoRequest;
import com.example.curs3projectback.dto.act.ReplacementActRequest;
import com.example.curs3projectback.dto.act.ReplacementActResponse;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.mapper.ReplacementActMapper;
import com.example.curs3projectback.model.MeterInfo;
import com.example.curs3projectback.model.ReplacementAct;
import com.example.curs3projectback.repository.AddressRepository;
import com.example.curs3projectback.repository.ReplacementActRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class ReplacementActService {

    private final ReplacementActRepository replacementActRepository;
    private final AddressRepository addressRepository;
    private final ReplacementActMapper replacementActMapper;
    private final CurrentUserService currentUserService;
    private final CurrentTenantService currentTenantService;
    private final TaskAccessService taskAccessService;

    @Transactional
    public ReplacementActResponse create(ReplacementActRequest request, Authentication authentication) {
        var user = currentUserService.getRequiredUser(authentication);
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Creating replacement act taskId={} tenant={} userId={}", request.getTaskId(), tenant.getCode(), user.getId());
        var task = taskAccessService.getAssignedTask(request.getTaskId(), user, tenant);
        var address = addressRepository.findByIdAndTenant_Id(request.getAddressId(), tenant.getId())
                .orElseThrow(() -> new ResourceNotFoundException("Address not found"));

        ReplacementAct act = ReplacementAct.builder()
                .task(task)
                .tenant(tenant)
                .address(address)
                .accountNumber(request.getAccountNumber())
                .installationDate(request.getInstallationDate())
                .oldMeter(mapInfo(request.getOldMeter()))
                .newMeter(mapInfo(request.getNewMeter()))
                .build();

        replacementActRepository.save(act);
        log.info("Replacement act created id={} tenant={}", act.getId(), tenant.getCode());
        return replacementActMapper.toResponse(act);
    }

    @Transactional(readOnly = true)
    public List<ReplacementActResponse> findMyActs(Authentication authentication) {
        var user = currentUserService.getRequiredUser(authentication);
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Loading replacement acts userId={} tenant={}", user.getId(), tenant.getCode());
        return replacementActMapper.toResponses(
                replacementActRepository.findAllByTask_AssigneeAndTenant_Id(user, tenant.getId())
        );
    }

    private MeterInfo mapInfo(MeterInfoRequest request) {
        if (request == null) {
            return null;
        }
        return MeterInfo.builder()
                .brand(request.getBrand())
                .serialNumber(request.getSerialNumber())
                .readings(request.getReadings())
                .build();
    }
}
