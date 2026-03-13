package com.example.curs3projectback.service;

import jakarta.annotation.PostConstruct;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.UUID;

@Slf4j
@Service
public class FileStorageService {

    private final Path uploadDir;

    public FileStorageService(@Value("${file.upload-dir}") String uploadDir) {
        this.uploadDir = Path.of(uploadDir);
    }

    @PostConstruct
    public void init() throws IOException {
        if (!Files.exists(uploadDir)) {
            Files.createDirectories(uploadDir);
        }
        log.info("File storage initialized at {}", uploadDir.toAbsolutePath());
    }

    public String store(MultipartFile file) {
        try {
            String filename = UUID.randomUUID() + "_" + file.getOriginalFilename();
            Path target = uploadDir.resolve(filename);
            Files.copy(file.getInputStream(), target, StandardCopyOption.REPLACE_EXISTING);
            log.info("Stored file {} at {}", file.getOriginalFilename(), target.toAbsolutePath());
            return filename;
        } catch (IOException e) {
            throw new RuntimeException("Failed to store file", e);
        }
    }

    public byte[] load(String filename) {
        try {
            log.info("Loading file {} from {}", filename, uploadDir.toAbsolutePath());
            return Files.readAllBytes(uploadDir.resolve(filename));
        } catch (IOException e) {
            throw new RuntimeException("Failed to load file", e);
        }
    }
}
