package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.photo.PhotoResponse;
import com.example.curs3projectback.dto.photo.PhotoUploadRequest;
import com.example.curs3projectback.service.PhotoService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/api/photos")
@RequiredArgsConstructor
public class PhotoController {

    private final PhotoService photoService;

    @PostMapping(consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public ResponseEntity<PhotoResponse> upload(@Valid @RequestPart("meta") PhotoUploadRequest request,
                                                @RequestPart("file") MultipartFile file,
                                                Authentication authentication) {
        log.info("HTTP POST /api/photos userId={} inspectionActId={} replacementActId={} fileName={}",
                authentication.getName(), request.getInspectionActId(), request.getReplacementActId(), file.getOriginalFilename());
        return ResponseEntity.ok(photoService.upload(request, file, authentication));
    }

    @GetMapping("/inspection/{actId}")
    public ResponseEntity<List<PhotoResponse>> listInspection(@PathVariable Long actId, Authentication authentication) {
        log.info("HTTP GET /api/photos/inspection/{} userId={}", actId, authentication.getName());
        return ResponseEntity.ok(photoService.listForInspection(actId, authentication));
    }

    @GetMapping("/replacement/{actId}")
    public ResponseEntity<List<PhotoResponse>> listReplacement(@PathVariable Long actId, Authentication authentication) {
        log.info("HTTP GET /api/photos/replacement/{} userId={}", actId, authentication.getName());
        return ResponseEntity.ok(photoService.listForReplacement(actId, authentication));
    }

    @GetMapping("/{photoId}")
    public ResponseEntity<byte[]> download(@PathVariable Long photoId, Authentication authentication) {
        log.info("HTTP GET /api/photos/{} userId={}", photoId, authentication.getName());
        var content = photoService.load(photoId, authentication);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "inline; filename=photo_" + photoId + ".jpg")
                .contentType(MediaType.IMAGE_JPEG)
                .body(content);
    }

    @DeleteMapping("/{photoId}")
    public ResponseEntity<Void> delete(@PathVariable Long photoId, Authentication authentication) {
        log.info("HTTP DELETE /api/photos/{} userId={}", photoId, authentication.getName());
        photoService.delete(photoId, authentication);
        return ResponseEntity.noContent().build();
    }
}

