package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.Photo;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface PhotoRepository extends JpaRepository<Photo, Long> {
    List<Photo> findAllByInspectionAct_Id(Long actId);
    List<Photo> findAllByReplacementAct_Id(Long actId);
}

