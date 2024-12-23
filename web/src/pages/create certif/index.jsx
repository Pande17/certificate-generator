import React, { useEffect, useState } from 'react';
import { useForm, Controller, useFieldArray } from 'react-hook-form';
import { Form, Input, DatePicker, Button, InputNumber, Select, message } from 'antd';
import MainLayout from '../MainLayout/Layout';
import { Sertifikat, Kompetensi, Signature } from '../api middleware';
import { useNavigate } from 'react-router-dom';

function MyForm() {
	const [data, setData] = useState([]);
	const [signatureData, setSignatureData] = useState([]);
	const [selectedSignature, setSelectedSignature] = useState(null);
	const [isSignatureSelected, setIsSignatureSelected] = useState(false);
	const [skkni, setSkkni] = useState('');
	const [divisi, setDivisi] = useState('');
	const { control, handleSubmit, reset, setValue } = useForm({
		defaultValues: {
			hardSkill: [],
			softSkill: [],
			selectedCompetenceId: '',
		},
	});

	const navigate = useNavigate();
	const { fields: hardSkillFields, replace: replaceHardSkill } = useFieldArray({
		control,
		name: 'hardSkill',
	});

	const { fields: softSkillFields, replace: replaceSoftSkill } = useFieldArray({
		control,
		name: 'softSkill',
	});

	const calculateTotalSkillScore = (hardSkills, softSkills) => {
		const totalHardSkillsScore = Array.isArray(hardSkills) ? hardSkills.reduce((acc, skill) => acc + (skill.skill_score || 0), 0) : 0; // Pastikan hardSkills adalah array, jika tidak, set default ke 0

		const totalSoftSkillsScore = Array.isArray(softSkills) ? softSkills.reduce((acc, skill) => acc + (skill.skill_score || 0), 0) : 0; // Pastikan softSkills adalah array, jika tidak, set default ke 0

		return totalHardSkillsScore + totalSoftSkillsScore;
	};

	const onSubmit = async (formData) => {
		const totalSkillScore = calculateTotalSkillScore(
			formData.hardSkill, // Pastikan ini adalah array
			formData.softSkill // Pastikan ini adalah array
		);

    try {
      const formattedData = {
        savedb: true,
        page_name: "page2a",
        zoom: 1.367,
        data: {
          sertif_name: formData.sertifikat,
          nama_peserta: formData.nama,
          kompeten_bidang: formData.fieldOfStudy,
          kompetensi: data.find(
            (item) => item._id === formData.selectedCompetenceId
          )?.nama_kompetensi,
          meet_time: formData.meetingTime,
          skkni: formData.skkni,
          validation: formData.validation,
          valid_date: {
            valid_start: formData.expiredTimeStart,
            valid_end: formData.expiredTimeEnd,
            valid_total: formData.validTime,
          },
          total_meet: formData.totalMeeting,
          kode_referral: {
            divisi: formData.codeReferralFieldOfStudy,
          },
          hard_skills: {
            skills: Array.isArray(formData.hardSkill)
              ? formData.hardSkill.map((skill) => ({
                  skill_name: skill.skill_name,
                  skill_jp: skill.jp,
                  skill_score: skill.skillScore,
                  description: skill.combined_units.split("\n").map((line) => {
                    const [unit_code, unit_title] = line.split(" - ");
                    return { unit_code, unit_title };
                  }),
                }))
              : [],
            total_skill_jp:
              formData.hardSkill?.reduce(
                (acc, skill) => acc + (skill.jp || 0),
                0
              ) || 0,
            total_skill_score: totalSkillScore, // Replace with actual computation if necessary
          },
          soft_skills: {
            skills: Array.isArray(formData.softSkill)
              ? formData.softSkill.map((skill) => ({
                  skill_name: skill.skill_name,
                  skill_jp: skill.jp,
                  skill_score: skill.skillScore,
                  description: skill.combined_units.split("\n").map((line) => {
                    const [unit_code, unit_title] = line.split(" - ");
                    return { unit_code, unit_title };
                  }),
                }))
              : [],
            total_skill_jp:
              formData.softSkill?.reduce(
                (acc, skill) => acc + (skill.jp || 0),
                0
              ) || 0,
            total_skill_score: totalSkillScore, // Replace with actual computation if necessary
          },
          signature: {
            config_name: formData.config_name,
            logo: formData.logo,
            role: formData.role,
            signature: formData.linkGambarPenandatangan,
            name: formData.namaPenandatangan,
            stamp: formData.stamp,
          },

          total_jp:
            (formData.hardSkill?.reduce(
              (acc, skill) => acc + (skill.jp || 0),
              0
            ) || 0) +
            (formData.softSkill?.reduce(
              (acc, skill) => acc + (skill.jp || 0),
              0
            ) || 0),
        },
      };

			const response = await Sertifikat.post('', formattedData);

			if (response.status === 200) {
				console.log(data);
				message.success('Certificate added successfully!');
				reset(); // Clear the form
			} else {
				console.log('Ada masalah dengan respons:', response);
			}
		} catch (error) {
			console.log(data);
			console.log('Error adding certificate:', error);
			message.error('Failed to add certificate. Please try again.');
		}

		navigate(`/dashboard`);
	};

	const { Option } = Select;

	useEffect(() => {
		// Fetch competence data
		const fetchCompetence = async () => {
			try {
				const response = await Kompetensi.get('/');
				setData(response.data.data);
			} catch (error) {
				console.log('Error fetching competence data:', error);
			}
		};

		// Fetch signature data
		const fetchSignature = async () => {
			try {
				const response = await Signature.get('/');
				setSignatureData(response.data.data);
			} catch (error) {
				console.log('Error fetching signature data:', error);
			}
		};

		fetchCompetence();
		fetchSignature();
	}, []);

	const fetchCompetence = async (competenceId) => {
		try {
			const response = await Kompetensi.get(`/${competenceId}`);

			const { hard_skills = [], soft_skills = [], skkni = '', divisi = '' } = response.data.data || {};

			const newHardSkills = hard_skills.map((hardSkill) => ({
				skill_name: hardSkill.skill_name || '',
				combined_units: hardSkill.description.map((unit) => `${unit.unit_code} - ${unit.unit_title}`).join('\n'),
			}));

			const newSoftSkills = soft_skills.map((softSkill) => ({
				skill_name: softSkill.skill_name || '',
				combined_units: softSkill.description.map((unit) => `${unit.unit_code} - ${unit.unit_title}`).join('\n'),
			}));

			replaceHardSkill(newHardSkills);
			replaceSoftSkill(newSoftSkills);

			// Simpan skkni dan divisi ke state
			setSkkni(skkni);
			setDivisi(divisi);
		} catch (err) {
			console.log(err);
		}
	};

	const handleCompetenceChange = (value) => {
		reset((prevValues) => ({
			...prevValues, // Pertahankan nilai sebelumnya
			selectedCompetenceId: value,
			hardSkill: [],
			softSkill: [],
		}));
		fetchCompetence(value);
	};

	const fetchSignatureId = async (SignatureId) => {
		try {
			const response = await Signature.get(`/${SignatureId}`);
			return response.data.data;
		} catch (error) {
			console.error('Error fetching signature:', error);
			return null;
		}
	};

	const handleSignatureChange = async (value) => {
		const signature = await fetchSignatureId(value);
		if (signature) {
			setSelectedSignature(signature);
			setIsSignatureSelected(true);

			setValue('configName', signature.config_name || '');
			setValue('namaPenandatangan', signature.name || '');
			setValue('role', signature.role || '');
			setValue('logo', signature.logo || '');
			setValue('linkGambarPenandatangan', signature.signature || '');
			setValue('stamp', signature.stamp || '');
		} else {
			setIsSignatureSelected(false);
		}
	};

  return (
    <MainLayout>
      <Form
        layout="vertical"
        style={{
          width: "95%",
          maxHeight: "100vh",
          overflowY: "scroll",
          backgroundColor: "white",
          padding: "40px",
          borderRadius: "20px",
          margin: "auto",
        }}
        onFinish={handleSubmit(onSubmit)}
      >
        <div className="text-center font-Poppins font-bold text-xl">
          Buat Sertifikat
        </div>
        <Form.Item label="Judul Sertifikat" required>
          <Controller
            name="sertifikat"
            control={control}
            rules={{ required: "Judul diperlukan" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan judul"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>
        <Form.Item label="Nama" required>
          <Controller
            name="nama"
            control={control}
            rules={{ required: "Nama diperlukan" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Made Rendy Putra Mahardika"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Bidang Studi" required>
          <Controller
            name="fieldOfStudy"
            control={control}
            rules={{ required: "Field of Study diperlukan" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan Bidang Studi"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Total Tahun" required>
          <Controller
            name="validTime"
            control={control}
            rules={{ required: "Valid Time diperlukan" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="2 Tahun"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

				<Form.Item label="Waktu Expired (Mulai)" required>
					<Controller name="expiredTimeStart" control={control} rules={{ required: 'Waktu mulai diperlukan' }} render={({ field }) => <Input {...field} placeholder="Pilih waktu mulai" style={{ width: '100%', height: '50px' }} />} />
				</Form.Item>

        <Form.Item label="Waktu Expired (Selesai)" required>
          <Controller
            name="expiredTimeEnd"
            control={control}
            rules={{ required: "Waktu selesai diperlukan" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Pilih waktu selesai"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>
        <Form.Item label="Total Pertemuan" required>
          <Controller
            name="totalMeeting"
            control={control}
            rules={{ required: "Total Pertemuan diperlukan" }}
            render={({ field }) => (
              <InputNumber
                {...field}
                placeholder="contoh: 14"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Total Waktu Pertemuan" required>
          <Controller
            name="meetingTime"
            control={control}
            rules={{ required: "Total Waktu Pertemuan diperlukan" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="contoh: 13"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>
        <Form.Item label="Waktu dan Tempat Pengesahan" required>
          <Controller
            name="validation"
            control={control}
            rules={{ required: "Validation diperlukan" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="contoh: 13"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <h1 className="text-center font-Poppins text-2xl font-medium p-6">
          Pilih Kompetensi
        </h1>
        <Form.Item required>
          <Controller
            name="selectedCompetenceId"
            control={control}
            render={({ field }) => (
              <Select
                placeholder="Pilih Kompetensi"
                {...field}
                style={{ width: "100%", height: "50px" }}
                onChange={(value) => {
                  field.onChange(value);
                  handleCompetenceChange(value);
                }}
              >
                <Option value="" disabled>
                  Pilih Kompetensi
                </Option>
                {data.map((competence) => (
                  <Option key={competence._id} value={competence._id}>
                    {competence.nama_kompetensi || ""}
                  </Option>
                ))}
              </Select>
            )}
          />
        </Form.Item>

				{skkni && (
					<Form.Item label="SKKNI">
						<Input value={skkni} readOnly style={{ width: '100%', height: '50px' }} />
					</Form.Item>
				)}

				{divisi && (
					<Form.Item label="Divisi">
						<Input value={divisi} readOnly style={{ width: '100%', height: '50px' }} />
					</Form.Item>
				)}

				{hardSkillFields.length > 0 && (
					<div>
						<h2 className="font-Poppins text-2xl font-medium text-center p-6">Hard Skills</h2>
						{hardSkillFields.map((skill, index) => (
							<div key={index} style={{ marginBottom: '20px' }}>
								<label>{`Hardskill ${index + 1}`}</label>

								{/* Skill Name Input */}
								<Controller
									name={`hardSkill[${index}].skill_name`}
									control={control}
									render={({ field }) => (
										<Input
											{...field}
											placeholder="Skill Name"
											readOnly
											style={{
												marginBottom: '10px',
												width: '100%',
												height: '50px',
											}}
										/>
									)}
								/>

								{/* Unit Code and Title Input */}
								<Controller
									name={`hardSkill[${index}].combined_units`}
									control={control}
									render={({ field }) => (
										<Input.TextArea
											{...field}
											rows={4}
											placeholder="Unit Code and Title"
											readOnly
											style={{
												marginBottom: '10px',
												width: '100%',
											}}
										/>
									)}
								/>

                {/* JP Input for each hard skill */}
                <Controller
                  name={`hardSkill[${index}].jp`}
                  control={control}
                  render={({ field }) => (
                    <InputNumber
                      {...field}
                      placeholder="JP per Skill"
                      style={{
                        width: "100%",
                        height: "50px",
                      }}
                    />
                  )}
                />
                <Controller
                  name={`hardSkill[${index}].skillScore`}
                  control={control}
                  render={({ field }) => (
                    <InputNumber
                      {...field}
                      placeholder="Score"
                      style={{
                        width: "100%",
                        height: "50px",
                      }}
                    />
                  )}
                />
              </div>
            ))}
          </div>
        )}

				{softSkillFields.length > 0 && (
					<div>
						<h2 className="font-Poppins text-2xl font-medium text-center p-6">Soft Skills</h2>
						{softSkillFields.map((skill, index) => (
							<div key={index} style={{ marginBottom: '20px' }}>
								<label>{`Softskill ${index + 1}`}</label>

								{/* Skill Name Input */}
								<Controller
									name={`softSkill[${index}].skill_name`}
									control={control}
									render={({ field }) => (
										<Input
											{...field}
											placeholder="Skill Name"
											readOnly
											style={{
												marginBottom: '10px',
												width: '100%',
												height: '50px',
											}}
										/>
									)}
								/>

								{/* Unit Code and Title Input */}
								<Controller
									name={`softSkill[${index}].combined_units`}
									control={control}
									render={({ field }) => (
										<Input.TextArea
											{...field}
											rows={4}
											placeholder="Unit Code and Title"
											readOnly
											style={{
												marginBottom: '10px',
												width: '100%',
											}}
										/>
									)}
								/>

                <Controller
                  name={`softSkill[${index}].jp`}
                  control={control}
                  render={({ field }) => (
                    <InputNumber
                      {...field}
                      placeholder="JP per Skill"
                      style={{
                        width: "100%",
                        height: "50px",
                      }}
                    />
                  )}
                />

								<Controller
									name={`softSkill[${index}].skillScore`}
									control={control}
									render={({ field }) => (
										<InputNumber
											{...field}
											placeholder="Score"
											style={{
												width: '100%',
												height: '50px',
											}}
										/>
									)}
								/>
							</div>
						))}
					</div>
				)}

				<Form.Item required>
					<Controller
						name="selectedSignatureId"
						control={control}
						render={({ field }) => (
							<Select {...field} placeholder="Pilih Template Paraf" onChange={handleSignatureChange} style={{ width: '100%', height: '50px' }}>
								<Option value="" disabled>
									Pilih Tanda Tangan
								</Option>
								{signatureData.map((signature) => (
									<Option key={signature._id} value={signature._id}>
										{signature.config_name}
									</Option>
								))}
							</Select>
						)}
					/>
				</Form.Item>


				{isSignatureSelected && selectedSignature && (
					<>
						<Form.Item label="Nama Display" required>
							<Controller name="configName" control={control} render={({ field }) => <Input {...field} readOnly style={{ width: '100%', height: '50px' }} />} />
						</Form.Item>

						<Form.Item label="Nama penandatangan" required>
							<Controller name="namaPenandatangan" control={control} render={({ field }) => <Input {...field} readOnly style={{ width: '100%', height: '50px' }} />} />
						</Form.Item>

            <Form.Item label="Jabatan Penandatangan" required>
              <Controller
                name="role"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    readOnly
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Stamp Perusahaan" required>
              <Controller
                name="stamp"
                control={control}
                render={({ field }) => (
                  <>
                    <Input
                      {...field}
                      readOnly
                      style={{ width: "100%", height: "50px" }}
                    />
                    {field.value && (
                      <div style={{ marginTop: "10px" }}>
                        <img
                          src={field.value}
                          alt="Stamp perusahaan"
                          style={{
                            width: "200px",
                            height: "200px",
                            border: "solid",
                            borderColor: "black",
                          }}
                        />
                      </div>
                    )}
                  </>
                )}
              />
            </Form.Item>

            <Form.Item label="Link logo" required>
              <Controller
                name="logo"
                control={control}
                render={({ field }) => (
                  <>
                    <Input
                      {...field}
                      readOnly
                      style={{ width: "100%", height: "50px" }}
                    />
                    {field.value && (
                      <div style={{ marginTop: "10px" }}>
                        <img
                          src={field.value}
                          alt="Logo perusahaan"
                          style={{
                            width: "200px",
                            height: "200px",
                            border: "solid",
                            borderColor: "black",
                          }}
                        />
                      </div>
                    )}
                  </>
                )}
              />
            </Form.Item>

            <Form.Item label="Link gambar penandatangan" required>
              <Controller
                name="linkGambarPenandatangan"
                control={control}
                render={({ field }) => (
                  <>
                    <Input
                      {...field}
                      readOnly
                      style={{ width: "100%", height: "50px" }}
                    />
                    {field.value && (
                      <div style={{ marginTop: "10px" }}>
                        <img
                          src={field.value}
                          alt="Tandatangan orang terkait"
                          style={{
                            width: "200px",
                            height: "200px",
                            border: "solid",
                            borderColor: "black",
                          }}
                        />
                      </div>
                    )}
                  </>
                )}
              />
            </Form.Item>
          </>
        )}

				<Form.Item>
					<Button type="primary" htmlType="submit">
						Submit
					</Button>
				</Form.Item>
			</Form>
		</MainLayout>
	);
}

export default MyForm;
