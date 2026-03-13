package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.act.ReplacementActRequest;
import com.example.curs3projectback.dto.act.ReplacementActResponse;
import com.example.curs3projectback.service.ReplacementActService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/api/acts/replacement")
@RequiredArgsConstructor
public class ReplacementActController {

    private final ReplacementActService replacementActService;

    @PostMapping
    public ResponseEntity<ReplacementActResponse> create(@Valid @RequestBody ReplacementActRequest request,
                                                         Authentication authentication) {
        log.info("HTTP POST /api/acts/replacement userId={} taskId={} addressId={}",
                authentication.getName(), request.getTaskId(), request.getAddressId());
        return ResponseEntity.ok(replacementActService.create(request, authentication));
    }

    @GetMapping
    public ResponseEntity<List<ReplacementActResponse>> myActs(Authentication authentication) {
        log.info("HTTP GET /api/acts/replacement userId={}", authentication.getName());
        return ResponseEntity.ok(replacementActService.findMyActs(authentication));
    }
}

