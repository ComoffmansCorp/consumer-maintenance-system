package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.act.MeterInfoRequest;
import com.example.curs3projectback.dto.act.ReplacementActRequest;
import com.example.curs3projectback.dto.act.ReplacementActResponse;
import com.example.curs3projectback.exception.BadRequestException;
import com.example.curs3projectback.exception.ResourceNotFoundException;
import com.example.curs3projectback.mapper.ReplacementActMapper;
import com.example.curs3projectback.model.MeterInfo;
import com.example.curs3projectback.model.ReplacementAct;
import com.example.curs3projectback.repository.AddressRepository;
import com.example.curs3projectback.repository.ReplacementActRepository;
import com.example.curs3projectback.repository.TaskRepository;
import com.example.curs3projectback.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@RequiredArgsConstructor
public class ReplacementActService {

    private final ReplacementActRepository replacementActRepository;
    private final TaskRepository taskRepository;
    private final AddressRepository addressRepository;
    private final UserRepository userRepository;
    private final ReplacementActMapper replacementActMapper;

    @Transactional
    public ReplacementActResponse create(ReplacementActRequest request, Authentication authentication) {
        var user = userRepository.findByUsername(authentication.getName())
                .orElseThrow(() -> new ResourceNotFoundException("Пользователь не найден"));
        var task = taskRepository.findById(request.getTaskId())
                .orElseThrow(() -> new ResourceNotFoundException("Задача не найдена"));
        if (!task.getAssignee().equals(user)) {
            throw new BadRequestException("Задача принадлежит другому исполнителю");
        }
        var address = addressRepository.findById(request.getAddressId())
                .orElseThrow(() -> new ResourceNotFoundException("Адрес не найден"));

        ReplacementAct act = ReplacementAct.builder()
                .task(task)
                .address(address)
                .accountNumber(request.getAccountNumber())
                .installationDate(request.getInstallationDate())
                .oldMeter(mapInfo(request.getOldMeter()))
                .newMeter(mapInfo(request.getNewMeter()))
                .build();

        replacementActRepository.save(act);
        return replacementActMapper.toResponse(act);
    }

    public List<ReplacementActResponse> findMyActs(Authentication authentication) {
        var user = userRepository.findByUsername(authentication.getName())
                .orElseThrow(() -> new ResourceNotFoundException("Пользователь не найден"));
        var acts = replacementActRepository.findAllByTask_Assignee(user);
        return replacementActMapper.toResponses(acts);
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

