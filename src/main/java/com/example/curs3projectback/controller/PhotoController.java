package com.example.curs3projectback.controller;

import com.example.curs3projectback.dto.photo.PhotoResponse;
import com.example.curs3projectback.dto.photo.PhotoUploadRequest;
import com.example.curs3projectback.service.PhotoService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

@RestController
@RequestMapping("/api/photos")
@RequiredArgsConstructor
public class PhotoController {

    private final PhotoService photoService;

    @PostMapping(consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public ResponseEntity<PhotoResponse> upload(@RequestPart("meta") PhotoUploadRequest request,
                                                @RequestPart("file") MultipartFile file) {
        return ResponseEntity.ok(photoService.upload(request, file));
    }

    @GetMapping("/inspection/{actId}")
    public ResponseEntity<List<PhotoResponse>> listInspection(@PathVariable Long actId) {
        return ResponseEntity.ok(photoService.listForInspection(actId));
    }

    @GetMapping("/replacement/{actId}")
    public ResponseEntity<List<PhotoResponse>> listReplacement(@PathVariable Long actId) {
        return ResponseEntity.ok(photoService.listForReplacement(actId));
    }

    @GetMapping("/{photoId}")
    public ResponseEntity<byte[]> download(@PathVariable Long photoId) {
        var content = photoService.load(photoId);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "inline; filename=photo_" + photoId + ".jpg")
                .contentType(MediaType.IMAGE_JPEG)
                .body(content);
    }

    @DeleteMapping("/{photoId}")
    public ResponseEntity<Void> delete(@PathVariable Long photoId) {
        photoService.delete(photoId);
        return ResponseEntity.noContent().build();
    }
}

