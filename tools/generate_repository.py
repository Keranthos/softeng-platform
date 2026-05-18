#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
代码生成工具 - Repository生成器
从数据库表结构自动生成Go语言的Repository接口和实现代码

使用方法:
    python tools/generate_repository.py <table_name> [--output-dir <dir>]

示例:
    python tools/generate_repository.py tools
    python tools/generate_repository.py projects --output-dir internal/repository
"""

import sys
import os
import re
import argparse
from typing import List, Dict, Tuple

# 数据库配置（从环境变量或配置文件读取）
DB_CONFIG = {
    'host': os.getenv('DB_HOST', 'localhost'),
    'port': int(os.getenv('DB_PORT', 3306)),
    'user': os.getenv('DB_USER', 'root'),
    'password': os.getenv('DB_PASSWORD', ''),
    'database': os.getenv('DB_NAME', 'softeng'),
}

# Go类型映射
GO_TYPE_MAP = {
    'INT': 'int',
    'INTEGER': 'int',
    'BIGINT': 'int64',
    'VARCHAR': 'string',
    'TEXT': 'string',
    'CHAR': 'string',
    'TIMESTAMP': 'time.Time',
    'DATETIME': 'time.Time',
    'DATE': 'time.Time',
    'BOOLEAN': 'bool',
    'BOOL': 'bool',
    'FLOAT': 'float64',
    'DOUBLE': 'float64',
    'DECIMAL': 'float64',
}

# 表名到实体名的映射规则
def table_to_entity(table_name: str) -> str:
    """将表名转换为实体名（首字母大写，单数形式）"""
    # 移除复数形式
    if table_name.endswith('s'):
        entity = table_name[:-1]
    else:
        entity = table_name
    
    # 首字母大写
    return entity.capitalize()

def get_primary_key(columns: List[Dict]) -> Tuple[str, str]:
    """获取主键字段名和类型"""
    for col in columns:
        if col.get('is_primary', False):
            col_type = col['type'].upper()
            go_type = GO_TYPE_MAP.get(col_type.split('(')[0], 'string')
            return col['name'], go_type
    # 默认主键
    return 'id', 'int'

def generate_repository_interface(table_name: str, columns: List[Dict], pk_name: str, pk_type: str) -> str:
    """生成Repository接口代码"""
    entity_name = table_to_entity(table_name)
    
    interface_code = f'''package repository

import (
	"context"
)

type {entity_name}Repository interface {{
	// GetByID 根据ID获取记录
	GetByID(ctx context.Context, {pk_name} {pk_type}) (map[string]interface{{}}, error)
	
	// Create 创建新记录
	Create(ctx context.Context, userID int, data map[string]interface{{}}) (map[string]interface{{}}, error)
	
	// Update 更新记录
	Update(ctx context.Context, {pk_name} {pk_type}, data map[string]interface{{}}) (map[string]interface{{}}, error)
	
	// Delete 删除记录
	Delete(ctx context.Context, {pk_name} {pk_type}) error
	
	// GetList 获取列表（带分页）
	GetList(ctx context.Context, limit, offset int) ([]map[string]interface{{}}, error)
}}
'''
    return interface_code

def generate_repository_impl(table_name: str, columns: List[Dict], pk_name: str, pk_type: str) -> str:
    """生成Repository实现代码"""
    entity_name = table_to_entity(table_name)
    struct_name = entity_name.lower() + 'Repository'
    
    # 生成SELECT字段列表
    select_fields = []
    scan_fields = []
    for col in columns:
        col_name = col['name']
        select_fields.append(f"t.{col_name}")
        # 处理可能为NULL的字段
        if col.get('nullable', True):
            scan_fields.append(f"var {col_name} sql.NullString")
        else:
            col_type = col['type'].upper().split('(')[0]
            go_type = GO_TYPE_MAP.get(col_type, 'string')
            scan_fields.append(f"var {col_name} {go_type}")
    
    # 生成INSERT字段
    insert_fields = [col['name'] for col in columns if col['name'] != pk_name and 'created_at' not in col['name']]
    insert_placeholders = ', '.join(['?' for _ in insert_fields])
    
    # 生成UPDATE字段
    update_fields = [col['name'] for col in columns if col['name'] != pk_name and 'created_at' not in col['name']]
    update_set = ', '.join([f"{f} = ?" for f in update_fields])
    
    # 生成扫描字段声明
    scan_declarations = []
    for col in columns:
        col_name = col['name']
        if col.get('nullable', True):
            scan_declarations.append(f"var {col_name} sql.NullString")
        else:
            col_type = col['type'].upper().split('(')[0]
            go_type = GO_TYPE_MAP.get(col_type, 'string')
            scan_declarations.append(f"var {col_name} {go_type}")
    
    # 生成字段列表
    field_names = [col['name'] for col in columns]
    
    # 生成字段映射
    field_mappings = []
    for col in columns:
        col_name = col['name']
        if col.get('nullable', True):
            field_mappings.append(f'\t\t"{col_name}": {col_name}.String,')
        else:
            field_mappings.append(f'\t\t"{col_name}": {col_name},')
    
    # 生成列表扫描字段声明
    list_scan_declarations = []
    for col in columns:
        col_name = col['name']
        col_type = col['type'].upper().split('(')[0]
        go_type = GO_TYPE_MAP.get(col_type, 'string')
        list_scan_declarations.append(f"\t\tvar {col_name} {go_type}")
    
    impl_code = f'''package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type {struct_name} struct {{
	db *Database
}}

func New{entity_name}Repository(db *Database) {entity_name}Repository {{
	return &{struct_name}{{db: db}}
}}

func (r *{struct_name}) GetByID(ctx context.Context, {pk_name} {pk_type}) (map[string]interface{{}}, error) {{
	query := `SELECT {', '.join(select_fields)} FROM {table_name} WHERE {pk_name} = ?`
	
	row := r.db.QueryRowContext(ctx, query, {pk_name})
	
	// 扫描字段
{chr(10).join([f'	{decl}' for decl in scan_declarations])}
	
	if err := row.Scan({', '.join(field_names)}); err != nil {{
		if err == sql.ErrNoRows {{
			return nil, fmt.Errorf("{entity_name} not found")
		}}
		return nil, fmt.Errorf("failed to get {entity_name.lower()}: %w", err)
	}}
	
	// 构建返回的map
	result := map[string]interface{{}}{{
		"{pk_name}": {pk_name},
{chr(10).join(field_mappings)}
	}}
	
	return result, nil
}}

func (r *{struct_name}) Create(ctx context.Context, userID int, data map[string]interface{{}}) (map[string]interface{{}}, error) {{
	query := `INSERT INTO {table_name} ({', '.join(insert_fields)}, created_at) VALUES ({insert_placeholders}, NOW())`
	
	args := []interface{{}}{{
{chr(10).join([f'		data["{field}"],' for field in insert_fields])}
	}}
	
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {{
		return nil, fmt.Errorf("failed to create {entity_name.lower()}: %w", err)
	}}
	
	id, err := result.LastInsertId()
	if err != nil {{
		return nil, fmt.Errorf("failed to get insert id: %w", err)
	}}
	
	// 返回创建的记录
	return r.GetByID(ctx, {pk_type}(id))
}}

func (r *{struct_name}) Update(ctx context.Context, {pk_name} {pk_type}, data map[string]interface{{}}) (map[string]interface{{}}, error) {{
	query := `UPDATE {table_name} SET {update_set}, updated_at = NOW() WHERE {pk_name} = ?`
	
	args := []interface{{}}{{
{chr(10).join([f'		data["{field}"],' for field in update_fields])}
		{pk_name},
	}}
	
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {{
		return nil, fmt.Errorf("failed to update {entity_name.lower()}: %w", err)
	}}
	
	// 返回更新后的记录
	return r.GetByID(ctx, {pk_name})
}}

func (r *{struct_name}) Delete(ctx context.Context, {pk_name} {pk_type}) error {{
	query := `DELETE FROM {table_name} WHERE {pk_name} = ?`
	
	_, err := r.db.ExecContext(ctx, query, {pk_name})
	if err != nil {{
		return fmt.Errorf("failed to delete {entity_name.lower()}: %w", err)
	}}
	
	return nil
}}

func (r *{struct_name}) GetList(ctx context.Context, limit, offset int) ([]map[string]interface{{}}, error) {{
	query := `SELECT {', '.join(select_fields)} FROM {table_name} ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {{
		return nil, fmt.Errorf("failed to query {entity_name.lower()} list: %w", err)
	}}
	defer rows.Close()
	
	var items []map[string]interface{{}}
	for rows.Next() {{
		// 扫描字段
{chr(10).join(list_scan_declarations)}
		
		if err := rows.Scan({', '.join(field_names)}); err != nil {{
			return nil, fmt.Errorf("failed to scan {entity_name.lower()}: %w", err)
		}}
		
		items = append(items, map[string]interface{{}}{{
{chr(10).join([f'\t\t\t"{col["name"]}": {col["name"]},' for col in columns])}
		}})
	}}
	
	return items, nil
}}
'''
    
    return impl_code

def parse_schema_file(schema_file: str, table_name: str) -> List[Dict]:
    """从schema.sql文件中解析表结构"""
    columns = []
    
    try:
        with open(schema_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # 查找表定义
        table_pattern = rf'CREATE TABLE.*?{table_name}.*?\((.*?)\)'
        match = re.search(table_pattern, content, re.DOTALL | re.IGNORECASE)
        
        if not match:
            print(f"警告: 在schema.sql中未找到表 {table_name} 的定义")
            return columns
        
        table_def = match.group(1)
        
        # 解析字段定义
        lines = table_def.split('\n')
        for line in lines:
            line = line.strip()
            if not line or line.startswith('--') or line.startswith('INDEX') or line.startswith('FOREIGN KEY') or line.startswith('UNIQUE KEY'):
                continue
            
            # 解析字段: name TYPE [NULL|NOT NULL] [COMMENT '...']
            field_match = re.match(r'(\w+)\s+(\w+(?:\([^)]+\))?)\s*(.*)', line, re.IGNORECASE)
            if field_match:
                col_name = field_match.group(1)
                col_type = field_match.group(2)
                col_attrs = field_match.group(3)
                
                # 检查是否为主键
                is_primary = 'PRIMARY KEY' in col_attrs.upper() or col_name == f'{table_name}_id'
                
                # 检查是否可空
                nullable = 'NOT NULL' not in col_attrs.upper()
                
                columns.append({
                    'name': col_name,
                    'type': col_type,
                    'is_primary': is_primary,
                    'nullable': nullable,
                })
    
    except FileNotFoundError:
        print(f"警告: 未找到schema文件 {schema_file}")
    
    return columns

def main():
    parser = argparse.ArgumentParser(description='生成Go Repository代码')
    parser.add_argument('table_name', help='数据库表名')
    parser.add_argument('--output-dir', default='internal/repository', help='输出目录')
    parser.add_argument('--schema-file', default='database/schema.sql', help='数据库schema文件路径')
    
    args = parser.parse_args()
    
    # 解析表结构
    print(f"正在解析表结构: {args.table_name}...")
    columns = parse_schema_file(args.schema_file, args.table_name)
    
    if not columns:
        print(f"错误: 无法解析表 {args.table_name} 的结构")
        print("提示: 请确保表名正确，或在schema.sql中定义了该表")
        sys.exit(1)
    
    print(f"找到 {len(columns)} 个字段")
    
    # 获取主键
    pk_name, pk_type = get_primary_key(columns)
    print(f"主键: {pk_name} ({pk_type})")
    
    # 生成代码
    entity_name = table_to_entity(args.table_name)
    print(f"生成 {entity_name}Repository...")
    
    interface_code = generate_repository_interface(args.table_name, columns, pk_name, pk_type)
    impl_code = generate_repository_impl(args.table_name, columns, pk_name, pk_type)
    
    # 确保输出目录存在
    os.makedirs(args.output_dir, exist_ok=True)
    
    # 写入文件
    output_file = os.path.join(args.output_dir, f'{args.table_name}_repository_generated.go')
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write("// Code generated by tools/generate_repository.py. DO NOT EDIT.\n")
        f.write("// 此文件由代码生成工具自动生成，请勿手动编辑\n\n")
        f.write(interface_code)
        f.write("\n\n")
        f.write(impl_code)
    
    print(f"✅ 代码已生成: {output_file}")
    print(f"\n提示: 生成的是基础代码，你可能需要:")
    print(f"  1. 根据业务需求调整方法签名")
    print(f"  2. 添加业务逻辑")
    print(f"  3. 处理特殊字段类型")
    print(f"  4. 添加错误处理和验证")

if __name__ == '__main__':
    main()

