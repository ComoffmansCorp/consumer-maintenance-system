package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.act.InspectionActRequest;
import com.example.curs3projectback.dto.act.InspectionActResponse;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.mapper.InspectionActMapper;
import com.example.curs3projectback.model.InspectionAct;
import com.example.curs3projectback.model.Meter;
import com.example.curs3projectback.repository.AddressRepository;
import com.example.curs3projectback.repository.InspectionActRepository;
import com.example.curs3projectback.repository.OrganizationRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class InspectionActService {

    private final InspectionActRepository inspectionActRepository;
    private final AddressRepository addressRepository;
    private final OrganizationRepository organizationRepository;
    private final InspectionActMapper inspectionActMapper;
    private final CurrentUserService currentUserService;
    private final CurrentTenantService currentTenantService;
    private final TaskAccessService taskAccessService;

    @Transactional
    public InspectionActResponse create(InspectionActRequest request, Authentication authentication) {
        var user = currentUserService.getRequiredUser(authentication);
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Creating inspection act taskId={} tenant={} userId={}", request.getTaskId(), tenant.getCode(), user.getId());
        var task = taskAccessService.getAssignedTask(request.getTaskId(), user, tenant);
        var address = addressRepository.findByIdAndTenant_Id(request.getAddressId(), tenant.getId())
                .orElseThrow(() -> new ResourceNotFoundException("Address not found"));
        var consumer = request.getConsumerId() == null ? null :
                organizationRepository.findByIdAndTenant_Id(request.getConsumerId(), tenant.getId())
                        .orElseThrow(() -> new ResourceNotFoundException("Organization not found"));

        InspectionAct act = InspectionAct.builder()
                .task(task)
                .tenant(tenant)
                .address(address)
                .consumer(consumer)
                .inspectionDate(request.getInspectionDate())
                .inspectionType(request.getInspectionType())
                .notes(request.getNotes())
                .build();

        List<Meter> meters = new ArrayList<>();
        if (request.getMeters() != null) {
            request.getMeters().forEach(meterRequest -> meters.add(Meter.builder()
                    .inspectionAct(act)
                    .type(meterRequest.getType())
                    .serialNumber(meterRequest.getSerialNumber())
                    .manufactureYear(meterRequest.getManufactureYear())
                    .verificationDate(meterRequest.getVerificationDate())
                    .sealState(meterRequest.getSealState())
                    .transformationRatio(meterRequest.getTransformationRatio())
                    .build()));
        }
        act.setMeters(meters);
        inspectionActRepository.save(act);
        log.info("Inspection act created id={} tenant={} meters={}", act.getId(), tenant.getCode(), meters.size());
        return inspectionActMapper.toResponse(act);
    }

    @Transactional(readOnly = true)
    public List<InspectionActResponse> findMyActs(Authentication authentication) {
        var user = currentUserService.getRequiredUser(authentication);
        var tenant = currentTenantService.getRequiredTenant(authentication);
        log.info("Loading inspection acts userId={} tenant={}", user.getId(), tenant.getCode());
        return inspectionActMapper.toResponses(
                inspectionActRepository.findAllByTask_AssigneeAndTenant_Id(user, tenant.getId())
        );
    }
}
