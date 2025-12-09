package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.act.InspectionActRequest;
import com.example.curs3projectback.dto.act.InspectionActResponse;
import com.example.curs3projectback.exception.BadRequestException;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.mapper.InspectionActMapper;
import com.example.curs3projectback.model.InspectionAct;
import com.example.curs3projectback.model.Meter;
import com.example.curs3projectback.repository.*;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.List;

@Service
@RequiredArgsConstructor
public class InspectionActService {

    private final InspectionActRepository inspectionActRepository;
    private final AddressRepository addressRepository;
    private final OrganizationRepository organizationRepository;
    private final TaskRepository taskRepository;
    private final UserRepository userRepository;
    private final InspectionActMapper inspectionActMapper;

    @Transactional
    public InspectionActResponse create(InspectionActRequest request, Authentication authentication) {
        var user = userRepository.findByUsername(authentication.getName())
                .orElseThrow(() -> new ResourceNotFoundException("Пользователь не найден"));
        var task = taskRepository.findById(request.getTaskId())
                .orElseThrow(() -> new ResourceNotFoundException("Задача не найдена"));
        if (!task.getAssignee().equals(user)) {
            throw new BadRequestException("Задача принадлежит другому исполнителю");
        }
        var address = addressRepository.findById(request.getAddressId())
                .orElseThrow(() -> new ResourceNotFoundException("Адрес не найден"));
        var consumer = request.getConsumerId() == null ? null :
                organizationRepository.findById(request.getConsumerId())
                        .orElseThrow(() -> new ResourceNotFoundException("Потребитель не найден"));

        InspectionAct act = InspectionAct.builder()
                .task(task)
                .address(address)
                .consumer(consumer)
                .inspectionDate(request.getInspectionDate())
                .inspectionType(request.getInspectionType())
                .notes(request.getNotes())
                .build();

        List<Meter> meters = new ArrayList<>();
        if (request.getMeters() != null) {
            request.getMeters().forEach(meterRequest -> {
                var meter = Meter.builder()
                        .inspectionAct(act)
                        .type(meterRequest.getType())
                        .serialNumber(meterRequest.getSerialNumber())
                        .manufactureYear(meterRequest.getManufactureYear())
                        .verificationDate(meterRequest.getVerificationDate())
                        .sealState(meterRequest.getSealState())
                        .transformationRatio(meterRequest.getTransformationRatio())
                        .build();
                meters.add(meter);
            });
        }
        act.setMeters(meters);
        inspectionActRepository.save(act);
        return inspectionActMapper.toResponse(act);
    }

    public List<InspectionActResponse> findMyActs(Authentication authentication) {
        var user = userRepository.findByUsername(authentication.getName())
                .orElseThrow(() -> new ResourceNotFoundException("Пользователь не найден"));
        var acts = inspectionActRepository.findAllByTask_Assignee(user);
        return inspectionActMapper.toResponses(acts);
    }
}

