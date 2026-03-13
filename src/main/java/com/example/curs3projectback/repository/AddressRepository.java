package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.Address;
import com.example.curs3projectback.model.Tenant;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface AddressRepository extends JpaRepository<Address, Long> {
    List<Address> findAllByStreetContainingIgnoreCaseOrderByStreet(String street);
    List<Address> findAllByTenantAndStreetContainingIgnoreCaseOrderByStreet(Tenant tenant, String street);
    Optional<Address> findByIdAndTenant_Id(Long id, Long tenantId);
}

