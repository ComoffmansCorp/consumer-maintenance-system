package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.Address;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface AddressRepository extends JpaRepository<Address, Long> {
    List<Address> findAllByStreetContainingIgnoreCaseOrderByStreet(String street);
}

