import React, { useState, useEffect } from 'react';
import { List, Button, Modal, Form, Input, message, Spin } from 'antd';
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useAuth } from '../hooks/useAuth';
import { calendarMuxApi, type CalendarMux, type CreateCalendarMuxRequest } from '../api';

const { TextArea } = Input;

const CalendarMuxList: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const [calendarMuxes, setCalendarMuxes] = useState<CalendarMux[]>([]);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const fetchCalendarMuxes = async () => {
    if (!isAuthenticated) return;

    setLoading(true);
    try {
      const response = await calendarMuxApi.list();
      setCalendarMuxes(response.calendar_muxes);
    } catch (error) {
      message.error('Failed to load calendar muxes');
      console.error('Error fetching calendar muxes:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCalendarMuxes();
  }, [isAuthenticated]);

  const handleDelete = async (id: number) => {
    try {
      await calendarMuxApi.delete(id);
      message.success('Calendar mux deleted successfully');
      fetchCalendarMuxes();
    } catch (error) {
      message.error('Failed to delete calendar mux');
      console.error('Error deleting calendar mux:', error);
    }
  };

  const handleCreate = async (values: CreateCalendarMuxRequest) => {
    try {
      await calendarMuxApi.create(values);
      message.success('Calendar mux created successfully');
      setIsModalOpen(false);
      form.resetFields();
      fetchCalendarMuxes();
    } catch (error) {
      message.error('Failed to create calendar mux');
      console.error('Error creating calendar mux:', error);
    }
  };

  if (!isAuthenticated) {
    return (
      <div style={{ textAlign: 'center', padding: '40px' }}>
        Please log in to view your calendar muxes
      </div>
    );
  }

  return (
    <div>
      <h2>Calendar Muxes</h2>

      <Spin spinning={loading}>
        <List
          dataSource={calendarMuxes}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button
                  type="text"
                  danger
                  icon={<DeleteOutlined />}
                  onClick={() => handleDelete(item.id)}
                >
                  Delete
                </Button>,
              ]}
            >
              <List.Item.Meta
                title={item.name}
                description={item.description}
              />
            </List.Item>
          )}
        />
      </Spin>

      <Button
        type="primary"
        icon={<PlusOutlined />}
        onClick={() => setIsModalOpen(true)}
        style={{ marginTop: '20px' }}
      >
        Create Calendar Mux
      </Button>

      <Modal
        title="Create Calendar Mux"
        open={isModalOpen}
        onCancel={() => {
          setIsModalOpen(false);
          form.resetFields();
        }}
        onOk={() => form.submit()}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreate}
        >
          <Form.Item
            name="name"
            label="Name"
            rules={[
              { required: true, message: 'Please enter a name' },
              { max: 200, message: 'Name must be at most 200 characters' },
            ]}
          >
            <Input placeholder="Enter calendar mux name" />
          </Form.Item>

          <Form.Item
            name="description"
            label="Description"
            rules={[
              { max: 1000, message: 'Description must be at most 1000 characters' },
            ]}
          >
            <TextArea
              rows={4}
              placeholder="Enter calendar mux description"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default CalendarMuxList;
