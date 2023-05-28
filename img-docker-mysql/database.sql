CREATE TABLE test.params_data (
    deploy_id INT,
    app_name VARCHAR(255),
    type_params VARCHAR(20),
    version_params VARCHAR(10),
    values_params JSON
);

INSERT INTO test.params_data (deploy_id, app_name, type_params, version_params, values_params)
VALUES (1, 'example', 'ci/release', '3.1.0', '{"key":"value"}');